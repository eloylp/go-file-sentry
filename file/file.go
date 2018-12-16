package file

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	file.calculateTime()
	file.calculateSum()
	file.calculateFQDN()
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

func (file *File) GetFQDN() string {
	return file.FQDN
}

func (file *File) calculateFQDN() {
	const sysDirNameSeparator = "-"
	const systemDirNameDatePart string = "20060102150405"
	fileDatePart := file.Time.Format(systemDirNameDatePart)
	parts := []string{file.GetSum(), fileDatePart}
	file.FQDN = strings.Join(parts, sysDirNameSeparator)
}

func (file *File) calculateTime() {
	readFile, err := os.Stat(file.Path)
	if err != nil {
		log.Fatal(err)
	}
	file.Time = readFile.ModTime()
}
