package file

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type File struct {
	Path string
	FQDN string
	Sum  string
	Time time.Time
}

func NewFile(path string) *File {
	file := File{Path: path}
	file.loadMetadata()
	return &file
}

func (file *File) loadMetadata() {
	file.calculateSum()
}

func (file *File) GetSum() string {
	return file.Sum
}

func (file *File) calculateSum() {
	targetFile, err := os.Open(file.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer targetFile.Close()
	sum := md5.New()
	if _, err = io.Copy(sum, targetFile); err != nil {
		log.Fatal(err)
	}
	sumInBytes := sum.Sum(nil)[:16]
	sumString := hex.EncodeToString(sumInBytes)
	file.Sum = sumString
}

func (file *File) GetData() []byte {

	dat, err := ioutil.ReadFile(file.Path)
	if err != nil {
		log.Fatal(err)
	}
	return dat
}

func (file *File) GetName() string {
	return filepath.Base(file.Path)
}
