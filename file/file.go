package file

import (
	"io/ioutil"
	"log"
	. "time"
)

type File struct {
	Path string
	FQDN string
	Sum  string
	Time Time
}

func (file File) GetData() []byte {

	dat, err := ioutil.ReadFile(file.Path)
	if err != nil {
		log.Fatal(err)
	}
	return dat
}
