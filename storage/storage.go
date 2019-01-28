package storage

import (
	"github.com/eloylp/go-file-sentry/file"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const diffExtension string = ".diff"

func EnsureSlot(rootPath string, storageUnitContent *StorageUnit) {
	destination := filepath.Join(rootPath, storageUnitContent.Name(), storageUnitContent.FileFQDN())
	err := os.MkdirAll(destination, 0755)
	failIfError(err)
}

func failIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func EntryContent(rootPath string, storageUnit *StorageUnit) {
	const defaultPerms os.FileMode = 0666

	fileName := storageUnit.FileName()
	containerName := storageUnit.Name()
	targetFilePath := filepath.Join(rootPath, containerName, storageUnit.FileFQDN(), fileName)
	err := ioutil.WriteFile(targetFilePath, storageUnit.FileData(), defaultPerms)
	failIfError(err)
	err = ioutil.WriteFile(targetFilePath+diffExtension, storageUnit.DiffContent(), defaultPerms)
	failIfError(err)
}

type VersionNotFound struct {
	err string
}

func (v VersionNotFound) Error() string {
	return v.err
}

func NewVersionNotFoundError(err string) *VersionNotFound {
	return &VersionNotFound{err: err}
}

func LatestVersion(rootPath string, scannedFile *file.File) (StorageUnit, error) {

	storageUnit := StorageUnit{
		file: scannedFile,
	}
	fContainerPath := filepath.Join(rootPath, storageUnit.Name())
	storedVersions, err := ioutil.ReadDir(fContainerPath)

	if storedVersions == nil {
		return storageUnit, NewVersionNotFoundError("There`s no previous storage units.")
	}
	lastVersion := calculateLastVersionReference(storedVersions)
	fullFilePath := filepath.Join(fContainerPath, lastVersion.Name(), scannedFile.Name())
	requestedFile := file.NewFile(fullFilePath)
	diffContent, err := ioutil.ReadFile(fullFilePath + diffExtension)
	failIfError(err)
	storageUnit = StorageUnit{
		file:        requestedFile,
		diffContent: diffContent,
	}
	return storageUnit, err
}

func calculateLastVersionReference(storedVersions []os.FileInfo) os.FileInfo {
	const containerPartsSeparator string = "-"
	const containerTimePartLayout = "20060102150405"

	var lastTime time.Time
	var lastVersion os.FileInfo
	for _, storedVersion := range storedVersions {
		name := storedVersion.Name()
		parts := strings.Split(name, containerPartsSeparator)
		timeStamp := parts[0]
		parsedTime, err := time.Parse(containerTimePartLayout, timeStamp)
		failIfError(err)
		if parsedTime.After(lastTime) {
			lastTime = parsedTime
			lastVersion = storedVersion
		}
	}
	return lastVersion
}
