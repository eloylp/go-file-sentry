package storage

import (
	"github.com/eloylp/go-file-sentry/file"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	const osPathSepNormalized = "_"
	containerName := strings.Replace(inputFile.Path, string(os.PathSeparator), osPathSepNormalized, -1)
	containerName = normalizeForFilePath(containerName)
	if strings.HasPrefix(containerName, osPathSepNormalized) {
		containerName = strings.TrimLeft(containerName, osPathSepNormalized)
	}
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
