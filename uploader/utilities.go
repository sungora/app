package uploader

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"os"
	"time"
)

const (
	characterMap string = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// randString Псевто-случайной строки символов указанной длинны из заданного набора символов
func randString(lenght int) string {
	var buf bytes.Buffer

	buf.Reset()
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < lenght; i++ {
		l := rand.Intn(len(characterMap))
		buf.WriteByte(characterMap[l])
	}
	return buf.String()
}

// checkFile Проверка файла
func checkFile(fileName string) fileInfo {
	var ret fileInfo
	var err error
	var fi os.FileInfo
	var file *os.File
	var tm time.Time
	var ic image.Config

	fi, err = os.Stat(fileName)
	if err == nil {
		tm = fi.ModTime()
		ret.ModTime = &tm
		ret.Size = fi.Size()

		file, err = os.Open(fileName)
		if err != nil {
			return ret
		}
		defer file.Close()

		ic, _, err = image.DecodeConfig(file)
		if err == nil {
			ret.PictureWidth = int64(ic.Width)
			ret.PictureHeight = int64(ic.Height)
			ret.IsPicture = true
		}
	}

	return ret
}
