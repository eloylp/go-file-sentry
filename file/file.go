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
	path string
	fqdn string
	sum  string
	time time.Time
}

func NewFile(path string) *File {
	file := File{path: path}
	file.LoadMetadata()
	return &file
}

func (file *File) LoadMetadata() {
	file.calcTime()
	file.calcSum()
	file.calcFQDN()
}

func (file *File) Sum() string {
	return file.sum
}

func (file *File) calcSum() {
	targetFile, err := os.Open(file.path)
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
	file.sum = sumString
}

func (file *File) Data() []byte {

	dat, err := ioutil.ReadFile(file.path)
	if err != nil {
		log.Fatal(err)
	}
	return dat
}

func (file *File) Name() string {
	return filepath.Base(file.path)
}

func (file *File) FQDN() string {
	return file.fqdn
}

func (file *File) calcFQDN() {
	const sysDirNameSeparator = "-"
	const systemDirNameDatePart = "20060102150405"
	fileDatePart := file.time.Format(systemDirNameDatePart)
	parts := []string{fileDatePart, file.Sum()}
	file.fqdn = strings.Join(parts, sysDirNameSeparator)
}

func (file *File) calcTime() {
	readFile, err := os.Stat(file.path)
	if err != nil {
		log.Fatal(err)
	}
	file.time = readFile.ModTime()
}

func (file *File) Path() string {
	return file.path
}
