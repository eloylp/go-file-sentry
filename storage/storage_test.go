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

	folderPath := _test.CreateTestStorageFolder()
	defer _test.CleanFolder(folderPath)
	filePath := _test.WriteFile(folderPath, "test.txt", "Content")
	sampleFile := file.NewFile(filePath)
	storageUnit := storage.NewStorageUnit([]byte{}, sampleFile)
	storage.EnsureSlot(folderPath, storageUnit)
	expectedFolderPath := filepath.Join(
		folderPath,
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

	folderPath := _test.CreateTestStorageFolder()
	defer _test.CleanFolder(folderPath)
	fileName := "fileA"
	fileContent := "Content A	"
	filePath := _test.WriteFile(folderPath, fileName, fileContent)

	sampleFile := file.NewFile(filePath)
	sampleFileDiffContent := []byte("Differential patch")
	storageUnit := storage.NewStorageUnit(sampleFileDiffContent, sampleFile)
	containerName := _test.Md5(filePath)
	containerFolder := filepath.Join(
		folderPath,
		containerName,
		sampleFile.FQDN())

	err := os.MkdirAll(containerFolder, 0755)
	_test.FailIfError(err)

	storage.EntryContent(folderPath, storageUnit)
	expectedFilePath := filepath.Join(containerFolder, fileName)

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

	folderPath := _test.CreateFixedTestFolder("TestFindLatestVersion")
	defer _test.CleanFolder(folderPath)
	samplePath := _test.WriteFile(folderPath, "fstab", "Content C")
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

	folderPath := _test.CreateTestStorageFolder()
	defer _test.CleanFolder(folderPath)
	samplePath := _test.WriteFile(folderPath, "fstab", "Content C")
	sampleFile := file.NewFile(samplePath)
	rootPath := _test.CreateTestStorageFolder()
	defer _test.CleanFolder(rootPath)
	_, err := storage.LatestVersion(rootPath, sampleFile)
	expectedErrorMessage := "There`s no previous storage units."
	if err.Error() != expectedErrorMessage {
		t.Fatalf("We are expecting error '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}
