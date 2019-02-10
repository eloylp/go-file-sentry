package main_test

import (
	"github.com/eloylp/go-file-sentry/_test"
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/sentry"
	"io/ioutil"
	"path"
	"testing"
	"time"
)

func TestStartSentry(t *testing.T) {

	testRootFolder := _test.CreateFixedTestStorageFolder("main_listening")
	testFileFolder := _test.CreateFixedTestStorageFolder("main_listening_files")

	file1 := _test.WriteFileToTestFolder(testFileFolder, "file1.txt", "A text file !")
	cfg := config.NewConfig(
		testRootFolder,
		[]string{file1})
	sentry.Start(cfg)

	time.Sleep(time.Duration(1 * time.Second))
	_test.AppendDataToTestFile(file1, "more content")
	time.Sleep(time.Duration(1 * time.Second))
	_test.AppendDataToTestFile(file1, "even more content")
	time.Sleep(time.Duration(1 * time.Second))
	infos, err := ioutil.ReadDir(path.Join(testRootFolder, _test.CalculateMd5(file1)))
	_test.FailIfError(err)
	versions := len(infos)

	_test.CleanTestStorageFolder(testRootFolder)
	_test.CleanTestStorageFolder(testFileFolder)
	if versions != 2 {
		t.Fatalf("Expected versions are 2 , got %d", versions)
	}
}
