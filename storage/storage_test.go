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

func createTestStorageFolder() (string, string) {
	const testRoot string = "tmp"
	const prefix string = "go_file_sentry_test_"
	now := time.Now().Format("20060102150405000000.000000")
	testFolder := filepath.Join(string(os.PathSeparator), testRoot, prefix+now)
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

	fileTime, err := time.Parse("2006-01-02 15:04:05", "2018-01-01 13:43:54")
	failIfError(err)

	sampleFile := file.File{}
	sampleFile.Sum = "587399e23181c0a8862b1c8c2a2225a6"
	sampleFile.Time = fileTime
	sampleFile.Path = "/etc/mysql/my.conf"
	sampleFile.FQDN = "587399e23181c0a8862b1c8c2a2225a6-20180101134354"

	testFolderPath, _ := createTestStorageFolder()
	defer cleanTestStorageFolder(testFolderPath)

	storageUnit := storage.StorageUnit{File: sampleFile}
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
	filePath := writeFileToTestFolder(testFolderPath, "fileA", "Content A	")
	fileTime, _ := time.Parse("2006-01-02 15:04:05", "2018-01-01 13:43:54")
	sampleFile := file.File{}
	sampleFile.Sum = "587399e23181c0a8862b1c8c2a2225a6"
	sampleFile.Time = fileTime
	sampleFile.Path = filePath
	sampleFile.FQDN = "587399e23181c0a8862b1c8c2a2225a6-20180101134354"
	sampleFileDiffContent := []byte("Differential patch")
	storageUnit := storage.StorageUnit{
		File:        sampleFile,
		DiffContent: sampleFileDiffContent,
	}
	containerFolder := filepath.Join(
		testFolderPath,
		"tmp_go_file_sentry_test_"+testFolderTime+"_fileA",
		sampleFile.FQDN)
	containerFolder = strings.Replace(containerFolder, ".", "_", -1)

	err := os.MkdirAll(containerFolder, 0755)
	failIfError(err)

	storage.AddEntryContent(testFolderPath, storageUnit)

	expectedFilePath := filepath.Join(containerFolder, "fileA")

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
