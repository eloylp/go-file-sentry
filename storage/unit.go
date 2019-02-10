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

func (s *StorageUnit) DiffContent() []byte {
	return s.diffContent
}

func (s *StorageUnit) File() *file.File {
	return s.file
}

func (s *StorageUnit) FileFQDN() string {
	return s.File().FQDN()
}

func (s *StorageUnit) FileName() string {
	return s.File().Name()
}

func (s *StorageUnit) FileData() []byte {
	return s.File().Data()
}

func (s *StorageUnit) GetFilePath() string {
	return s.File().Path()
}

func (s *StorageUnit) Name() string {
	hasher := md5.New()
	hasher.Write([]byte(s.GetFilePath()))
	container := hex.EncodeToString(hasher.Sum(nil))
	return container
}
