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

func (unit *StorageUnit) DiffContent() []byte {
	return unit.diffContent
}

func (unit *StorageUnit) File() *file.File {
	return unit.file
}

func (unit *StorageUnit) FileFQDN() string {
	return unit.File().FQDN()
}

func (unit *StorageUnit) FileName() string {
	return unit.File().Name()
}

func (unit *StorageUnit) FileData() []byte {
	return unit.File().Data()
}

func (unit *StorageUnit) GetFilePath() string {
	return unit.File().Path()
}

func (unit *StorageUnit) Name() string {
	hasher := md5.New()
	hasher.Write([]byte(unit.GetFilePath()))
	containerName := hex.EncodeToString(hasher.Sum(nil))
	return containerName
}
