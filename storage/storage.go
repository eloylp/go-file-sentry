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

func AddNewEntry(rootPath string, storageUnitContent *StorageUnit) {
	destination := filepath.Join(rootPath, storageUnitContent.CalculateName(), storageUnitContent.GetFileFQDN())
	err := os.MkdirAll(destination, 0755)
	failIfError(err)
}

func failIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func AddEntryContent(rootPath string, storageUnit *StorageUnit) {
	const defaultPerms os.FileMode = 0666

	fileName := storageUnit.GetFileName()
	containerName := storageUnit.CalculateName()
	targetFilePath := filepath.Join(rootPath, containerName, storageUnit.GetFileFQDN(), fileName)
	err := ioutil.WriteFile(targetFilePath, storageUnit.GetFileData(), defaultPerms)
	failIfError(err)
	err = ioutil.WriteFile(targetFilePath+diffExtension, storageUnit.GetDiffContent(), defaultPerms)
	failIfError(err)
}

type VersionNotFound struct {
	err string
}

func NewVersionNotFoundError(err string) *VersionNotFound {
	return &VersionNotFound{err: err}
}

func (v VersionNotFound) Error() string {
	return v.err
}

func FindLatestVersion(rootPath string, scannedFile *file.File) (storageUnit StorageUnit, err error) {

	scannedFileStorageUnit := StorageUnit{
		file: scannedFile,
	}
	fileContainerDirAbsolutePath := filepath.Join(rootPath, scannedFileStorageUnit.CalculateName())
	storedVersions, err := ioutil.ReadDir(fileContainerDirAbsolutePath)

	if storedVersions == nil {
		return storageUnit, NewVersionNotFoundError("Theres no previous storage units.")
	}
	lastVersion := calculateLastVersionReference(storedVersions)
	fullFilePath := filepath.Join(fileContainerDirAbsolutePath, lastVersion.Name(), scannedFile.GetName())
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
		entryName := storedVersion.Name()
		parts := strings.Split(entryName, containerPartsSeparator)
		timeStampPart := parts[1]
		parsedTime, err := time.Parse(containerTimePartLayout, timeStampPart)
		failIfError(err)
		if parsedTime.After(lastTime) {
			lastTime = parsedTime
			lastVersion = storedVersion
		}
	}
	return lastVersion
}
