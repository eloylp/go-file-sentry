package version_test

import (
	"github.com/eloylp/go-file-sentry/_test"
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/version"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestMakeNewVersion(t *testing.T) {

	testFolder := _test.CreateFixedTestFolder("version_make_new")
	defer _test.CleanFolder(testFolder)
	testStorageFolder := _test.CreateFixedTestFolder("version_make_new_storage")
	defer _test.CleanFolder(testStorageFolder)

	expectedFileContent := "this.is.the.config=true"
	testConfFilePath := _test.WriteFile(testFolder, "daemon.config", expectedFileContent)
	testConfFile := file.NewFile(testConfFilePath)
	version.NewVersion(testStorageFolder, testConfFile)

	expectedContainerPath := calculateExpectedContainerPath(testStorageFolder, testConfFile)

	expectedFilePath := filepath.Join(expectedContainerPath, testConfFile.Name())
	assertFileContent(expectedFilePath, expectedFileContent, t)
	expectedFileDiffPath := expectedFilePath + ".diff"
	expectedFileDiffContent := []byte("")
	fileDiffContent, err := ioutil.ReadFile(expectedFileDiffPath)
	_test.FailIfError(err)
	assertDiffs(fileDiffContent, expectedFileDiffContent, t)

}

func assertFileContent(expectedFilePath string, expectedFileContent string, t *testing.T) {
	fileContent, err := ioutil.ReadFile(expectedFilePath)
	_test.FailIfError(err)
	if expectedFileContent != string(fileContent) {
		t.Fatalf("Expected file content is %s , got %s", expectedFileContent, fileContent)
	}
}

func assertDiffs(fileDiffContent []byte, expectedFileDiffContent []byte, t *testing.T) {
	if string(fileDiffContent) != string(expectedFileDiffContent) {
		t.Fatalf("Expected diff content '%s' does not match '%s'", expectedFileDiffContent, fileDiffContent)
	}
}

func calculateExpectedContainerPath(testStorageFolder string, testConfFile *file.File) string {
	return filepath.Join(
		testStorageFolder,
		_test.Md5(testConfFile.Path()),
		testConfFile.FQDN(),
	)
}

func TestMakeSecondVersion(t *testing.T) {

	testFolder := _test.CreateFixedTestFolder("version_make_new_2")
	defer _test.CleanFolder(testFolder)
	testStorageFolder := _test.CreateFixedTestFolder("version_make_new_storage_2")
	defer _test.CleanFolder(testStorageFolder)

	testConfFilePath := _test.WriteFile(testFolder, "daemon.config", "this.is.the.config=true")
	testConfFile := file.NewFile(testConfFilePath)
	version.NewVersion(testStorageFolder, testConfFile)
	expectedFileContent := "this.is.the.config=false"
	testConfFile2Path := _test.WriteFile(testFolder, "daemon.config", expectedFileContent)
	testConfFile2 := file.NewFile(testConfFile2Path)

	version.NewVersion(testStorageFolder, testConfFile2)
	expectedContainerPath := calculateExpectedContainerPath(testStorageFolder, testConfFile2)
	expectedFilePath := filepath.Join(expectedContainerPath, testConfFile.Name())

	assertFileContent(expectedFilePath, expectedFileContent, t)
	expectedFileDiffPath := filepath.Join(expectedContainerPath, testConfFile2.Name()+".diff")
	expectedFileDiffContent, err := ioutil.ReadFile(_test.GetTestResource("expectedDiff.diff"))
	_test.FailIfError(err)
	fileDiffContent, err := ioutil.ReadFile(expectedFileDiffPath)
	_test.FailIfError(err)
	assertDiffs(fileDiffContent, expectedFileDiffContent, t)
}
