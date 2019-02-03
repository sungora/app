package tool

import (
	"hash/crc64"
	"io/ioutil"
	"os"
	"path/filepath"
)

type controlFS struct {
	Files       map[string]uint64
	FilesChange map[string]uint64
	Path        string
	Ext         string
}

// NewControlFS Функция конструтор контроля изменений файлов по указанному пути
// Можно укзать фильтр по файлам (расширения файлов)
// Учитывается добавление, изменение и удаление файлов
func NewControlFS(path string, ext string) *controlFS {
	self := new(controlFS)
	self.Files = make(map[string]uint64)
	self.FilesChange = make(map[string]uint64)
	self.Path = path
	self.Ext = ext
	return self
}

// CheckSumMd5 проверка добавления, изменения и удаление файлов
func (self *controlFS) CheckSumMd5() (isChange bool, err error) {
	self.FilesChange = make(map[string]uint64)
	err = filepath.Walk(self.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// фильтр по расширению файла
		if self.Ext != "" && self.Ext != "*" && self.Ext != filepath.Ext(path) {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		// получаем контрольную сумму
		self.FilesChange[path] = crc64.Checksum(data, crc64.MakeTable(crc64.ECMA))
		// self.FilesChange[path] = fmt.Sprintf("%x", md5.Sum(data))
		// новые файлы
		if _, ok := self.Files[path]; ok == false {
			self.Files[path] = self.FilesChange[path]
			self.FilesChange[path] = 2
			isChange = true
			return nil
		}
		// существующие файлы
		if self.Files[path] != self.FilesChange[path] {
			self.Files[path] = self.FilesChange[path]
			self.FilesChange[path] = 1
			isChange = true
			return nil
		}
		return nil
	})
	if err != nil {
		return
	}
	// удаленные файлы
	for k, _ := range self.Files {
		if _, ok := self.FilesChange[k]; ok == false {
			delete(self.Files, k)
			self.FilesChange[k] = 0
			isChange = true
		}
	}
	return
}

// func main() {
// 	var path = "/home/konstantin/go/src/accounter/config/accountercron.toml"
// 	file, err := os.Open(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()
// 	hash := md5.New()
// 	_, err = io.Copy(hash, file)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%s MD5 checksum is %x \n", file.Name(), hash.Sum(nil))
// }
