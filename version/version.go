package version

import (
	"github.com/eloylp/go-file-sentry/diff"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/storage"
	"log"
)

func MakeNewVersion(rootPath string, file *file.File) {
	var diffFromPrevious string
	previousUnit, err := storage.FindLatestVersion(rootPath, file)
	switch err.(type) {
	case *storage.VersionNotFound:
		diffFromPrevious = ""
	case nil:
		diffFromPrevious = diff.GetDiffOfFiles(file, previousUnit.GetFile())
	default:
		log.Fatal(err)
	}

	unit := storage.NewStorageUnit([]byte(diffFromPrevious), file)
	storage.AddNewEntry(rootPath, unit)
	storage.AddEntryContent(rootPath, unit)
}
