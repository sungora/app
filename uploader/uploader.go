package uploader

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type uploader struct {
	Folder   string              // Папка хранения временных файлов
	Duration time.Duration       // Время жизни файлов
	Location *time.Location      // Текущая временная зона
	Keys     map[int64]time.Time // Все ключи загруженных файлов
}

const (
	tmpFolder string = `uploaderTemp` // Имя папки внутри временной папки
	tmpReload int    = 60 * 60        // Количество итераций после которых происходит перечитывание временной папки
)

var self *uploader

// Upload Загружает файл во временное хранилище и возвращает
// уникальный идентификатор файла
// По истечении времени жизни, если файл не будет востребован,
// он удаляется из временного хранилища автоматически
// Возвращаемые значения
// key - уникальный идентификатор файла
// err - ошибка загрузки или причины не работоспособности
var Upload func(req *http.Request, field string) (key string, err error)

// Get Получение ранее загруженного файла по ключу
var Get func(key string) (*Response, error)

// Delete Удаление всех данных по ключу
var Delete func(key string) error

// Инициализация
func init() {
	if self == nil {
		self = new(uploader)
	}
	self.Keys = make(map[int64]time.Time)

	Upload = func(req *http.Request, field string) (key string, err error) {
		key, err = self.upload(req, field)
		return key, err
	}

	Get = func(key string) (ret *Response, err error) {
		ret, err = self.get(key)
		return ret, err
	}

	Delete = func(key string) error {
		return self.del(key)
	}

	go self.cleaner()
}

// Init Инициализация библиотеки
// folder - путь к временной папке с файлами
// hour - время жизни загруженых файлов в минутах
// tl - информация о текущей временной зоне
func Init(folder string, minute int) (err error) {
	if folder == "" {
		return errors.New("Не корректно указана папка загрузки файлов")
	}
	if minute <= 0 {
		return errors.New("Не верно указан период жизни файла")
	}
	folder = folder + `/` + tmpFolder
	if err = os.MkdirAll(folder, 0777); err != nil {
		return
	}
	self.Folder = folder
	self.Duration = time.Minute * time.Duration(minute)
	self.Location = time.UTC
	return
}

// loadMultipartData Получение формы с файлами, enctype="multipart/form-data"
// Все файлы сохраняются в файловую систему без использования памяти
func loadMultipartData(req *http.Request, in *request) (*Response, error) {
	var ret *Response
	var err error
	var reader *multipart.Reader
	var part *multipart.Part
	var buf []byte

	if in == nil {
		err = errors.New("Не корректное обращение, не заполнен Request")
		return ret, err
	}

	// Начинаем чтение секций
	reader, err = req.MultipartReader()
	if err != nil {
		err = errors.New("Ошибка чтения данных: " + err.Error())
		return ret, err
	}

	// Критических ошибок больше не предвидится, можно создать ответ
	ret = new(Response)
	//	ret.DataMap = make(map[string]*data)

	// Обходим все данные передаваемые браузером
	for {
		var item *data = new(data)

		part, err = reader.NextPart()
		if err == io.EOF {
			break
		}

		// Если выборочное получение полей
		if in.Field != "" {
			var found bool
			// Имена полей формы могут быть на русском!
			if strings.EqualFold(part.FormName(), in.Field) == true {
				found = true
				ret.Field = in.Field
			}
			if found == false {
				continue
			}
		}

		// Инфа о форме
		item.Name = part.FormName()
		item.NameOriginal = part.FileName()
		item.ContentType = part.Header.Get("Content-Type")

		// Какая-то не предвиденная ошибка
		if err != nil {
			item.Error = err
			continue
		}

		// Если на входе ждёт файл
		if item.NameOriginal != "" || item.ContentType != "" {
			var file *os.File
			var rex *regexp.Regexp
			var info fileInfo
			var tmp []string

			item.NameFilePrefix = in.NameFilePrefix
			item.NameFile = item.NameFilePrefix + randString(16)

			// Получение расширения файла
			rex, err = regexp.Compile(`.*(\..+)$`)
			tmp = rex.FindStringSubmatch(item.NameOriginal)
			if len(tmp) > 1 {
				item.Extension = tmp[1]
			} else {
				item.Extension = `.data`
			}
			item.NameFile += item.Extension

			// Формирование куда сохранять и откуда запрашивать
			item.PathSys = in.PathSys + string(os.PathSeparator) + item.NameFile
			item.PathWeb = in.PathWeb + `/` + item.NameFile

			file, err = os.OpenFile(item.PathSys, os.O_CREATE|os.O_WRONLY, 0660)
			if err != nil {
				item.Error = err
				continue
			}
			defer file.Close()

			// Копирование файла из потока на диск без использования чтения в память
			item.Size, item.Error = io.Copy(file, part)
			if item.Error != nil {
				file.Close()
				continue
			}
			file.Close()

			// Проверка файла
			info = checkFile(item.PathSys)
			item.IsPicture = info.IsPicture
			item.Size = info.Size
			item.PictureWidth = info.PictureWidth
			item.PictureHeight = info.PictureHeight
		}

		// Если на входе обычное поле формы
		if item.NameOriginal == "" && item.ContentType == "" {
			buf, err = ioutil.ReadAll(part)
			item.Error = err
			item.Size = int64(len(buf))
			//item.Data.Write(buf)
		}

		ret.Data = *item
		//		ret.Data = append(ret.Data, item)
		//		ret.DataMap[item.Name] = item
	}
	if err == io.EOF {
		err = nil
	}
	return ret, err
}

// upload Загружает файл во временное хранилище и возвращает
// уникальный идентификатор файла
// По истечении времени жизни, если файл не будет востребован,
// он удаляется из временного хранилища автоматически
// Возвращаемые значения
// key - уникальный идентификатор файла
// err - ошибка загрузки или причины не работоспособности
func (self *uploader) upload(req *http.Request, field string) (key string, ret error) {
	var rq *request = new(request)
	var data *Response
	var err error
	var ok bool = false
	var ki int64

	ki = time.Now().In(self.Location).UnixNano()
	key = strconv.FormatInt(ki, 10)
	rq.PathSys = self.Folder
	rq.NameFilePrefix = key + `-`
	rq.Field = field

	data, err = loadMultipartData(req, rq)
	if err == nil {
		if data.Data.Size > 0 {
			ok = true

			// Сохранение структуры
			var file *os.File
			var buflen uint64 = 1 << 30 // 1Gb
			var buffer *bytes.Buffer
			var gb *gob.Encoder
			var fn string

			buffer = bytes.NewBuffer(make([]byte, 0, buflen))
			gb = gob.NewEncoder(buffer)
			err = gb.Encode(data)
			if err != nil {
				ret = errors.New("Не удалось сохранить структуру данных в бинарный формат: " + err.Error())
				os.Remove(data.Data.PathSys)
				key = ""
			} else {
				fn = strings.Replace(data.Data.PathSys, data.Data.NameFile, key+`.data`, -1)
				file, err = os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0660)
				if err != nil {
					ret = errors.New("Ошибка создания файла " + fn + ": " + err.Error())
					key = ""
					return key, ret
				}
				file.Write(buffer.Bytes())
				file.Close()
				self.Keys[ki] = time.Unix(0, ki).In(self.Location)
			}
		}
	}
	if ok == false {
		key = ""
		ret = errors.New("Не корректно передан файл, либо файл нулевой длинны")
	}
	return key, ret
}

// get Получение ранее загруженного файла по ключу
func (self *uploader) get(key string) (*Response, error) {
	var fnd string
	var fi os.FileInfo
	var fl *os.File
	var ret *Response
	var err error
	var gb *gob.Decoder

	// Формирование имени служебного файла
	fnd = self.Folder + string(os.PathSeparator) + key + `.data`
	// Проверка наличия файла
	fi, err = os.Stat(fnd)
	if err != nil || fi == nil {
		err = errors.New("Ключ не действительный")
		return ret, err
	}

	// Открытие файла
	fl, err = os.Open(fnd)
	if err != nil {
		err = errors.New("Ключ не действительный")
		return ret, err
	}
	defer fl.Close()

	var buflen uint64 = 1 << 30 // 1Gb
	var buf []byte
	var buffer *bytes.Buffer

	buffer = bytes.NewBuffer(make([]byte, 0, buflen))
	buf, err = ioutil.ReadAll(fl)
	buffer.Write(buf)

	// Чтение и декодирование бинарного сохранения в структуру
	ret = new(Response)
	gb = gob.NewDecoder(buffer)
	err = gb.Decode(ret)

	return ret, err
}

// del Удаление всех данных по ключу
func (self *uploader) del(key string) error {
	var ret error
	var fnd string
	var err error
	var resp *Response
	var ki int64

	ki, err = strconv.ParseInt(key, 0, 64)

	// Формирование имени служебного файла
	fnd = self.Folder + string(os.PathSeparator) + key + `.data`

	// Запрос данных по ключу
	resp, err = self.get(key)
	if err != nil {
		ret = err
	} else {
		ret = os.Remove(fnd)
		ret = os.Remove(resp.Data.PathSys)
	}
	if ret == nil {
		delete(self.Keys, ki)
	}

	return ret
}

// cleaner Очистка всех старых файлов
func (self *uploader) cleaner() {
	var di []os.FileInfo
	var err error
	var rex *regexp.Regexp
	var kn []string
	var ki int64
	var count int

	rex, err = regexp.Compile(`^(\d+)[\.\-]`)
	for {
		if self.Folder != "" && self.Duration > 0 && self.Location != nil {

			// Считываем ключи из папки только раз в несколько минут
			if count == 0 {
				di, err = ioutil.ReadDir(self.Folder)
				if err == nil {
					// Все ключи что есть в папке проверяются в карте
					for i := range di {
						kn = rex.FindStringSubmatch(di[i].Name())
						if len(kn) == 2 {
							ki, err = strconv.ParseInt(kn[1], 0, 64)
							if err == nil {
								if _, ok := self.Keys[ki]; ok == false {
									self.Keys[ki] = time.Unix(0, ki).In(self.Location)
								}

							}
						}
					}
					// Проверка наличия в карте не действительных ключей
					for k := range self.Keys {
						var found bool = false
						for i := range di {
							kn = rex.FindStringSubmatch(di[i].Name())
							if len(kn) == 2 {
								if kn[1] == strconv.FormatInt(k, 10) {
									found = true
								}
							}
						}
						if found == false {
							delete(self.Keys, k)
						}
					}
				}
			}

			// Проверка времени жизни и отправка на удаление
			for k := range self.Keys {
				if time.Now().In(self.Location).Sub(self.Keys[k]) > self.Duration {
					self.del(strconv.FormatInt(k, 10))
				}
			}
		}

		count++
		if count > tmpReload {
			count = 0
		}
		time.Sleep(time.Second)
	}
}
