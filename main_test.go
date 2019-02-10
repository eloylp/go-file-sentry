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

	root := _test.CreateFixedTestFolder("main_listening")
	fileFolder := _test.CreateFixedTestFolder("main_listening_files")

	file1 := _test.WriteFile(fileFolder, "file1.txt", "A text file !")
	cfg := config.NewConfig(
		root,
		[]string{file1})
	sentry.Start(cfg)

	time.Sleep(time.Duration(1 * time.Second))
	_test.AppendData(file1, "more content")
	time.Sleep(time.Duration(1 * time.Second))
	_test.AppendData(file1, "even more content")
	time.Sleep(time.Duration(1 * time.Second))

	infos, err := ioutil.ReadDir(path.Join(root, _test.Md5(file1)))
	_test.FailIfError(err)
	versions := len(infos)

	_test.CleanFolder(root)
	_test.CleanFolder(fileFolder)
	if versions != 2 {
		t.Fatalf("Expected versions are 2 , got %d", versions)
	}
}
