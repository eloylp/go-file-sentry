package storage

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/eloylp/go-file-sentry/file"
)

type StorageUnit struct {
	DiffContent []byte
	File        file.File
}

func (unit *StorageUnit) GetDiffContent() []byte {
	return unit.DiffContent
}

func (unit *StorageUnit) GetFileFQDN() string {
	return unit.File.GetFQDN()
}

func (unit *StorageUnit) GetFileName() string {
	return unit.File.GetName()
}

func (unit *StorageUnit) GetFileData() []byte {
	return unit.File.GetData()
}

func (unit *StorageUnit) GetFilePath() string {
	return unit.File.GetPath()
}

func (unit *StorageUnit) CalculateName() string {
	hasher := md5.New()
	hasher.Write([]byte(unit.GetFilePath()))
	containerName := hex.EncodeToString(hasher.Sum(nil))
	return containerName
}
