package storage_test

import (
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/storage"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeFileToTestFolder(testFolderPath string, name string, content string) string {
	filePath := filepath.Join(testFolderPath, string(os.PathSeparator), name)
	testFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	failIfError(err)
	defer testFile.Close()
	_, err = testFile.WriteString(content)
	failIfError(err)
	return filePath
}
func failIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const testRootFolderName string = "tmp"
const testFolderPrefix string = "go_file_sentry_test_"

func createTestStorageFolder() (string, string) {
	now := time.Now().Format("20060102150405000000.000000")
	testFolder := filepath.Join(string(os.PathSeparator), testRootFolderName, testFolderPrefix+now)
	os.Mkdir(testFolder, 0755)
	return testFolder, now
}

func cleanTestStorageFolder(path string) {
	err := os.RemoveAll(path)
	failIfError(err)
}

func fsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func TestAddNewEntry(t *testing.T) {

	sampleFile := file.File{
		Path: "/etc/mysql/my.conf",
		FQDN: "587399e23181c0a8862b1c8c2a2225a6-20180101134354",
	}
	testFolderPath, _ := createTestStorageFolder()
	defer cleanTestStorageFolder(testFolderPath)
	storageUnit := storage.StorageUnit{
		File: sampleFile,
	}
	storage.AddNewEntry(testFolderPath, storageUnit)
	expectedFolderPath := filepath.Join(testFolderPath, "etc_mysql_my_conf", "587399e23181c0a8862b1c8c2a2225a6-20180101134354")
	exist, err := fsExists(expectedFolderPath)

	if err != nil {
		t.Errorf(err.Error())
	}
	if !exist {
		t.Errorf("Cannot ensure that entry folder was added %s", expectedFolderPath)
	}
}

func TestAddEntryContent(t *testing.T) {

	testFolderPath, testFolderTime := createTestStorageFolder()
	defer cleanTestStorageFolder(testFolderPath)
	testFileName := "fileA"
	testFileContent := "Content A	"
	filePath := writeFileToTestFolder(testFolderPath, testFileName, testFileContent)
	sampleFile := file.File{
		Path: filePath,
		FQDN: "587399e23181c0a8862b1c8c2a2225a6-20180101134354",
	}
	sampleFileDiffContent := []byte("Differential patch")
	storageUnit := storage.StorageUnit{
		File:        sampleFile,
		DiffContent: sampleFileDiffContent,
	}
	const separator = "_"
	containerName := testRootFolderName + separator + testFolderPrefix + testFolderTime + separator + testFileName
	containerFolder := filepath.Join(
		testFolderPath,
		containerName,
		sampleFile.FQDN)
	containerFolder = strings.Replace(containerFolder, ".", separator, -1)

	err := os.MkdirAll(containerFolder, 0755)
	failIfError(err)

	storage.AddEntryContent(testFolderPath, storageUnit)
	expectedFilePath := filepath.Join(containerFolder, testFileName)

	exist, _ := fsExists(expectedFilePath)
	if !exist {
		t.Fatal("File doesnt exist in storage engine.")
	}

	expectedFileDiffPath := expectedFilePath + ".diff"
	exist, err = fsExists(expectedFileDiffPath)
	failIfError(err)
	if !exist {
		t.Fatal("File diff doesnt exist in storage engine.")
	}

	diffContent, err := ioutil.ReadFile(expectedFileDiffPath)
	failIfError(err)
	if string(diffContent) != string(sampleFileDiffContent) {
		t.Fatal("Unexpected file diff content.")
	}
}
