package file

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

type File struct {
	Path string
	FQDN string
	Sum  string
	Time time.Time
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
