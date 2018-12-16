package storage_test

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/storage"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
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

func createTestStorageFolder() string {
	uuidGen, _ := uuid.NewV4()
	testFolder := filepath.Join(string(os.PathSeparator), testRootFolderName, testFolderPrefix+uuidGen.String())
	os.Mkdir(testFolder, 0755)
	return testFolder
}

func cleanTestStorageFolder(path string) {
	err := os.RemoveAll(path)
	failIfError(err)
}

func calculateMd5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	containerName := hex.EncodeToString(hasher.Sum(nil))
	return containerName
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

func getTestResource(resourceName string) string {
	const testResourceDir = "test"
	return filepath.Join(testResourceDir, resourceName)
}

func TestAddNewEntry(t *testing.T) {

	sampleFile := file.File{
		Path: "/etc/mysql/my.conf",
		FQDN: "587399e23181c0a8862b1c8c2a2225a6-20180101134354",
	}
	testFolderPath := createTestStorageFolder()
	defer cleanTestStorageFolder(testFolderPath)
	storageUnit := storage.StorageUnit{
		File: sampleFile,
	}
	storage.AddNewEntry(testFolderPath, storageUnit)
	expectedFolderPath := filepath.Join(
		testFolderPath,
		"07ceeeacdce3f57f0c8164eb8ee21bec",
		"587399e23181c0a8862b1c8c2a2225a6-20180101134354",
	)
	exist, err := fsExists(expectedFolderPath)

	if err != nil {
		t.Errorf(err.Error())
	}
	if !exist {
		t.Errorf("Cannot ensure that entry folder was added %s", expectedFolderPath)
	}
}

func TestAddEntryContent(t *testing.T) {

	testFolderPath := createTestStorageFolder()
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
	containerName := calculateMd5(filePath)
	containerFolder := filepath.Join(
		testFolderPath,
		containerName,
		sampleFile.FQDN)

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

func TestFindLatestVersion(t *testing.T) {

	testFolderPath := getTestResource("root_sample")
	currentDir, _ := os.Getwd()
	testRootPath := filepath.Join(string(os.PathSeparator), currentDir, testFolderPath)
	// path := filepath.Join(testRootPath, "816a3ef77b62ade41c3f936c958aa555", "7b8fdf40404049204ed4feb3c8e99480-20180101134356", "fstab")
	path := "/etc/fstab"

	recoveredFile, err := storage.FindLatestVersion(testRootPath, file.File{
		Path: path,
		FQDN: "587399e23181c0a8862b1c8c2a2225a6-20180101134356",
	})
	failIfError(err)

	if "Content C" != string(recoveredFile.File.GetData()) {
		t.Fatal("Retrieved file is not the latest.")
	}

	if "+fake diff C" != string(recoveredFile.DiffContent) {
		t.Fatal("Retrieved file diff is not the latest.")
	}
}
