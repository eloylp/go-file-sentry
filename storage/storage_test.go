package storage_test

import (
	"github.com/eloylp/go-file-sentry/_test"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/storage"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureSlot(t *testing.T) {

	testFolderPath := _test.CreateTestStorageFolder()
	defer _test.CleanTestStorageFolder(testFolderPath)
	testFilePath := _test.WriteFileToTestFolder(testFolderPath, "test.txt", "Content")
	sampleFile := file.NewFile(testFilePath)
	storageUnit := storage.NewStorageUnit([]byte{}, sampleFile)
	storage.EnsureSlot(testFolderPath, storageUnit)
	expectedFolderPath := filepath.Join(
		testFolderPath,
		storageUnit.Name(),
		sampleFile.FQDN(),
	)
	exist, err := _test.FsExists(expectedFolderPath)

	if err != nil {
		t.Errorf(err.Error())
	}
	if !exist {
		t.Errorf("Cannot ensure that entry folder was added %s", expectedFolderPath)
	}
}

func TestEntryContent(t *testing.T) {

	testFolderPath := _test.CreateTestStorageFolder()
	defer _test.CleanTestStorageFolder(testFolderPath)
	testFileName := "fileA"
	testFileContent := "Content A	"
	filePath := _test.WriteFileToTestFolder(testFolderPath, testFileName, testFileContent)

	sampleFile := file.NewFile(filePath)
	sampleFileDiffContent := []byte("Differential patch")
	storageUnit := storage.NewStorageUnit(sampleFileDiffContent, sampleFile)
	containerName := _test.CalculateMd5(filePath)
	containerFolder := filepath.Join(
		testFolderPath,
		containerName,
		sampleFile.FQDN())

	err := os.MkdirAll(containerFolder, 0755)
	_test.FailIfError(err)

	storage.EntryContent(testFolderPath, storageUnit)
	expectedFilePath := filepath.Join(containerFolder, testFileName)

	exist, _ := _test.FsExists(expectedFilePath)
	if !exist {
		t.Fatal("file doesnt exist in storage engine.")
	}

	expectedFileDiffPath := expectedFilePath + ".diff"
	exist, err = _test.FsExists(expectedFileDiffPath)
	_test.FailIfError(err)
	if !exist {
		t.Fatal("file diff doesnt exist in storage engine.")
	}

	diffContent, err := ioutil.ReadFile(expectedFileDiffPath)
	_test.FailIfError(err)
	if string(diffContent) != string(sampleFileDiffContent) {
		t.Fatal("Unexpected file diff content.")
	}
}

func TestLatestVersion(t *testing.T) {

	testFolderPath := _test.CreateFixedTestStorageFolder("TestFindLatestVersion")
	defer _test.CleanTestStorageFolder(testFolderPath)
	samplePath := _test.WriteFileToTestFolder(testFolderPath, "fstab", "Content C")
	sampleFile := file.NewFile(samplePath)
	rootPath := _test.GetTestResource("root_sample")
	recoveredFile, err := storage.LatestVersion(rootPath, sampleFile)
	_test.FailIfError(err)

	if "Content C" != string(recoveredFile.FileData()) {
		t.Fatal("Retrieved file is not the latest.")
	}

	if "+fake diff C" != string(recoveredFile.DiffContent()) {
		t.Fatal("Retrieved file diff is not the latest.")
	}
}

func TestLatestVersionNotFound(t *testing.T) {

	testFolderPath := _test.CreateTestStorageFolder()
	defer _test.CleanTestStorageFolder(testFolderPath)
	samplePath := _test.WriteFileToTestFolder(testFolderPath, "fstab", "Content C")
	sampleFile := file.NewFile(samplePath)
	rootPath := _test.CreateTestStorageFolder()
	defer _test.CleanTestStorageFolder(rootPath)
	_, err := storage.LatestVersion(rootPath, sampleFile)
	expectedErrorMessage := "There`s no previous storage units."
	if err.Error() != expectedErrorMessage {
		t.Fatalf("We are expecting error '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}
