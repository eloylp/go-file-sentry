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
	time2 "time"
)

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

func normalizeForFilePath(input string) string {
	return strings.Replace(input, ".", "_", -1)
}

func AddEntryContent(rootPath string, storageUnit StorageUnit) {
	_, fileName := filepath.Split(storageUnit.File.Path)
	contName := calculateFileContainerName(storageUnit.File)
	targetFilePath := filepath.Join(rootPath, contName, storageUnit.File.FQDN, fileName)
	targetFilePath = normalizeForFilePath(targetFilePath)
	err := ioutil.WriteFile(targetFilePath, storageUnit.File.GetData(), 0666)
	failIfError(err)
	err = ioutil.WriteFile(targetFilePath+".diff", storageUnit.DiffContent, 0666)
	failIfError(err)
}

func FindLatestVersion(rootPath string, scannedFile file.File) (storageUnit StorageUnit, err error) {
	scannedFileDir := filepath.Join(rootPath, calculateFileContainerName(scannedFile))
	storedFiles, err := ioutil.ReadDir(scannedFileDir)
	failIfError(err)
	var lastTime time2.Time
	var lastFile os.FileInfo
	for _, storedFile := range storedFiles {
		entryName := storedFile.Name()
		sep := strings.Split(entryName, "-")
		parsedTime, err := time2.Parse("20060102150405", sep[1])
		failIfError(err)
		if parsedTime.After(lastTime) {
			lastTime = parsedTime
			lastFile = storedFile
		}
	}

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
