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
	file.calculateTime()
	file.calculateSum()
	file.calculateFQDN()
}

func (file *File) GetSum() string {
	return file.sum
}

func (file *File) calculateSum() {
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

func (file *File) GetData() []byte {

	dat, err := ioutil.ReadFile(file.path)
	if err != nil {
		log.Fatal(err)
	}
	return dat
}

func (file *File) GetName() string {
	return filepath.Base(file.path)
}

func (file *File) GetFQDN() string {
	return file.fqdn
}

func (file *File) calculateFQDN() {
	const sysDirNameSeparator = "-"
	const systemDirNameDatePart string = "20060102150405"
	fileDatePart := file.time.Format(systemDirNameDatePart)
	parts := []string{file.GetSum(), fileDatePart}
	file.fqdn = strings.Join(parts, sysDirNameSeparator)
}

func (file *File) calculateTime() {
	readFile, err := os.Stat(file.path)
	if err != nil {
		log.Fatal(err)
	}
	file.time = readFile.ModTime()
}

func (file *File) GetPath() string {
	return file.path
}
