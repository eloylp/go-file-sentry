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
	_ = os.Mkdir(testFolder, 0755)
	return testFolder
}

func createFixedTestStorageFolder(suffix string) string {
	testFolder := filepath.Join(string(os.PathSeparator), testRootFolderName, testFolderPrefix+suffix)
	_ = os.Mkdir(testFolder, 0755)
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

	testFolderPath := createTestStorageFolder()
	defer cleanTestStorageFolder(testFolderPath)
	testFilePath := writeFileToTestFolder(testFolderPath, "test.txt", "Content")
	sampleFile := file.NewFile(testFilePath)
	storageUnit := storage.StorageUnit{
		File: *sampleFile,
	}
	storage.AddNewEntry(testFolderPath, storageUnit)
	expectedFolderPath := filepath.Join(
		testFolderPath,
		storageUnit.CalculateName(),
		sampleFile.GetFQDN(),
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

	sampleFile := file.NewFile(filePath)
	sampleFileDiffContent := []byte("Differential patch")
	storageUnit := storage.StorageUnit{
		File:        *sampleFile,
		DiffContent: sampleFileDiffContent,
	}
	containerName := calculateMd5(filePath)
	containerFolder := filepath.Join(
		testFolderPath,
		containerName,
		sampleFile.GetFQDN())

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

	testFolderPath := createFixedTestStorageFolder("TestFindLatestVersion")
	samplePath := writeFileToTestFolder(testFolderPath, "fstab", "Content C")
	sampleFile := file.NewFile(samplePath)
	rootPath := getTestResource("root_sample")
	recoveredFile, err := storage.FindLatestVersion(rootPath, *sampleFile)
	failIfError(err)

	if "Content C" != string(recoveredFile.GetFileData()) {
		t.Fatal("Retrieved file is not the latest.")
	}

	if "+fake diff C" != string(recoveredFile.GetDiffContent()) {
		t.Fatal("Retrieved file diff is not the latest.")
	}
}
