package storage

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/scan"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const diffExtension string = ".diff"

type StorageUnit struct {
	DiffContent []byte
	File        file.File
}

func (unit *StorageUnit) GetDiffContent() []byte {
	return unit.DiffContent
}

func (unit *StorageUnit) GetFileFQDN() string {
	return unit.File.FQDN
}

func (unit *StorageUnit) GetFileName() string {
	return unit.File.GetName()
}

func (unit *StorageUnit) GetFileData() []byte {
	return unit.File.GetData()
}

func (unit *StorageUnit) GetFilePath() string {
	return unit.File.Path
}

func (unit *StorageUnit) CalculateName() string {
	hasher := md5.New()
	hasher.Write([]byte(unit.GetFilePath()))
	containerName := hex.EncodeToString(hasher.Sum(nil))
	return containerName
}

func AddNewEntry(rootPath string, storageUnitContent StorageUnit) {
	destination := filepath.Join(rootPath, storageUnitContent.CalculateName(), storageUnitContent.GetFileFQDN())
	err := os.MkdirAll(destination, 0755)
	failIfError(err)
}

func failIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func AddEntryContent(rootPath string, storageUnit StorageUnit) {
	const defaultPerms os.FileMode = 0666

	fileName := storageUnit.GetFileName()
	containerName := storageUnit.CalculateName()
	targetFilePath := filepath.Join(rootPath, containerName, storageUnit.GetFileFQDN(), fileName)
	err := ioutil.WriteFile(targetFilePath, storageUnit.GetFileData(), defaultPerms)
	failIfError(err)
	err = ioutil.WriteFile(targetFilePath+diffExtension, storageUnit.GetDiffContent(), defaultPerms)
	failIfError(err)
}

func FindLatestVersion(rootPath string, scannedFile file.File) (storageUnit StorageUnit, err error) {

	scannedFileStorageUnit := StorageUnit{
		File: scannedFile,
	}
	fileContainerDirAbsolutePath := filepath.Join(rootPath, scannedFileStorageUnit.CalculateName())
	storedVersions, err := ioutil.ReadDir(fileContainerDirAbsolutePath)
	failIfError(err)

	lastFile := calculateLastFileReference(storedVersions)

	fullFilePath := filepath.Join(fileContainerDirAbsolutePath, lastFile.Name(), scannedFile.GetName())
	requestedFile := scan.ScanFile(fullFilePath)
	diffContent, err := ioutil.ReadFile(fullFilePath + diffExtension)
	failIfError(err)
	storageUnit = StorageUnit{
		File:        requestedFile,
		DiffContent: diffContent,
	}
	return storageUnit, err
}

func calculateLastFileReference(storedVersions []os.FileInfo) os.FileInfo {
	const containerPartsSeparator string = "-"
	const containerTimePartLayout = "20060102150405"

	var lastTime time.Time
	var lastFile os.FileInfo
	for _, storedVersion := range storedVersions {
		entryName := storedVersion.Name()
		parts := strings.Split(entryName, containerPartsSeparator)
		timeStampPart := parts[1]
		parsedTime, err := time.Parse(containerTimePartLayout, timeStampPart)
		failIfError(err)
		if parsedTime.After(lastTime) {
			lastTime = parsedTime
			lastFile = storedVersion
		}
	}
	return lastFile
}
