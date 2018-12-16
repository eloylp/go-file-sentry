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
const defaultPerms os.FileMode = 0666

type StorageUnit struct {
	DiffContent []byte
	File        file.File
}

func AddNewEntry(rootPath string, storageUnitContent StorageUnit) {
	destination := filepath.Join(rootPath, calculateFileContainerName(storageUnitContent.File), storageUnitContent.File.FQDN)
	err := os.MkdirAll(destination, 0755)
	failIfError(err)
}

func failIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func calculateFileContainerName(inputFile file.File) string {
	hasher := md5.New()
	hasher.Write([]byte(inputFile.Path))
	containerName := hex.EncodeToString(hasher.Sum(nil))
	return containerName
}

func AddEntryContent(rootPath string, storageUnit StorageUnit) {
	fileName := storageUnit.File.GetName()
	containerName := calculateFileContainerName(storageUnit.File)
	targetFilePath := filepath.Join(rootPath, containerName, storageUnit.File.FQDN, fileName)

	err := ioutil.WriteFile(targetFilePath, storageUnit.File.GetData(), defaultPerms)
	failIfError(err)
	err = ioutil.WriteFile(targetFilePath+diffExtension, storageUnit.DiffContent, defaultPerms)
	failIfError(err)
}

func FindLatestVersion(rootPath string, scannedFile file.File) (storageUnit StorageUnit, err error) {
	scannedFileDir := filepath.Join(rootPath, calculateFileContainerName(scannedFile))
	storedFiles, err := ioutil.ReadDir(scannedFileDir)
	failIfError(err)
	lastFile := getLastFileByDate(storedFiles)

	fullFilePath := filepath.Join(scannedFileDir, lastFile.Name(), scannedFile.GetName())
	requestedFile := scan.ScanFile(fullFilePath)
	diffContent, err := ioutil.ReadFile(fullFilePath + ".diff")
	failIfError(err)
	storageUnit = StorageUnit{
		File:        requestedFile,
		DiffContent: diffContent,
	}
	return storageUnit, err
}

func getLastFileByDate(storedFileContainers []os.FileInfo) os.FileInfo {
	var lastTime time.Time
	var lastFile os.FileInfo
	for _, storedFileContainer := range storedFileContainers {
		entryName := storedFileContainer.Name()
		sep := strings.Split(entryName, "-")
		parsedTime, err := time.Parse("20060102150405", sep[1])
		failIfError(err)
		if parsedTime.After(lastTime) {
			lastTime = parsedTime
			lastFile = storedFileContainer
		}
	}
	return lastFile
}
