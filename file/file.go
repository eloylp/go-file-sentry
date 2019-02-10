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

func (f *File) LoadMetadata() {
	f.calcTime()
	f.calcSum()
	f.calcFQDN()
}

func (f *File) Sum() string {
	return f.sum
}

func (f *File) calcSum() {
	file, err := os.Open(f.path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	sum := md5.New()
	if _, err = io.Copy(sum, file); err != nil {
		log.Fatal(err)
	}
	sumData := sum.Sum(nil)[:16]
	sumString := hex.EncodeToString(sumData)
	f.sum = sumString
}

func (f *File) Data() []byte {

	dat, err := ioutil.ReadFile(f.path)
	if err != nil {
		log.Fatal(err)
	}
	return dat
}

func (f *File) Name() string {
	return filepath.Base(f.path)
}

func (f *File) FQDN() string {
	return f.fqdn
}

func (f *File) calcFQDN() {
	const sysDirNameSeparator = "-"
	const systemDirNameDatePart = "20060102150405"
	fileDatePart := f.time.Format(systemDirNameDatePart)
	parts := []string{fileDatePart, f.Sum()}
	f.fqdn = strings.Join(parts, sysDirNameSeparator)
}

func (f *File) calcTime() {
	readFile, err := os.Stat(f.path)
	if err != nil {
		log.Fatal(err)
	}
	f.time = readFile.ModTime()
}

func (f *File) Path() string {
	return f.path
}
