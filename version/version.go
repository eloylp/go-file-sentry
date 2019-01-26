package version

import (
	"github.com/eloylp/go-file-sentry/diff"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/storage"
	"log"
)

func NewVersion(rootPath string, file *file.File) {
	var diffFromPrevious string
	previousUnit, err := storage.LatestVersion(rootPath, file)
	switch err.(type) {
	case *storage.VersionNotFound:
		diffFromPrevious = ""
	case nil:
		diffFromPrevious = diff.DiffOfFiles(file, previousUnit.File())
	default:
		log.Fatal(err)
	}

	unit := storage.NewStorageUnit([]byte(diffFromPrevious), file)
	storage.EnsureSlot(rootPath, unit)
	storage.EntryContent(rootPath, unit)
}
