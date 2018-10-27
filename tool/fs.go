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
}

func NewControlFS() *controlFS {
	self := new(controlFS)
	self.Files = make(map[string]uint64)
	self.FilesChange = make(map[string]uint64)
	return self
}

// CheckSumMd5 контроль изменений файлов по указанному пути
// Учитывается добавление, изменение и удаление файлов
// Можно укзать фильтр по файлам (расширения файлов)
func (self *controlFS) CheckSumMd5(root string, ext string) (isChange bool, err error) {
	self.FilesChange = make(map[string]uint64)
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// фильтр по расширению файла
		if ext != "" && ext != "*" && ext != filepath.Ext(path) {
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
