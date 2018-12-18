package storage

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/eloylp/go-file-sentry/file"
)

type StorageUnit struct {
	diffContent []byte
	file        *file.File
}

func NewStorageUnit(diffContent []byte, file *file.File) *StorageUnit {
	return &StorageUnit{diffContent: diffContent, file: file}
}

func (unit *StorageUnit) GetDiffContent() []byte {
	return unit.diffContent
}

func (unit *StorageUnit) GetFile() *file.File {
	return unit.file
}

func (unit *StorageUnit) GetFileFQDN() string {
	return unit.GetFile().GetFQDN()
}

func (unit *StorageUnit) GetFileName() string {
	return unit.GetFile().GetName()
}

func (unit *StorageUnit) GetFileData() []byte {
	return unit.GetFile().GetData()
}

func (unit *StorageUnit) GetFilePath() string {
	return unit.GetFile().GetPath()
}

func (unit *StorageUnit) CalculateName() string {
	hasher := md5.New()
	hasher.Write([]byte(unit.GetFilePath()))
	containerName := hex.EncodeToString(hasher.Sum(nil))
	return containerName
}
