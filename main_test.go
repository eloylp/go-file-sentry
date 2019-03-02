package main_test

import (
	"github.com/eloylp/go-file-sentry/_test"
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/factory"
	"io/ioutil"
	"path"
	"sync"
	"testing"
	"time"
)

func TestStartSentry(t *testing.T) {
	/// TODO , NEEDS TO BE REWRITTEN
	t.Skip()
	root := _test.CreateFixedTestFolder("main_listening")
	fileFolder := _test.CreateFixedTestFolder("main_listening_files")

	file1 := _test.WriteFile(fileFolder, "file1.txt", "A text file !")
	cfg := config.NewConfig(
		root,
		[]string{file1},
		"")

	wg := new(sync.WaitGroup)
	factory.Watchers(cfg, wg)

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

func BenchmarkStartSentry(b *testing.B) {

	/// TODO , NEEDS TO BE REWRITTEN
	b.Skip()

	root := _test.CreateFixedTestFolder("main_listening")
	fileFolder := _test.CreateFixedTestFolder("main_listening_files")
	file := _test.WriteFile(fileFolder, "file1.txt", "A text file !")

	cfg := config.NewConfig(
		root,
		[]string{file}, "")
	wg := new(sync.WaitGroup)

	factory.Watchers(cfg, wg)
	time.Sleep(time.Duration(3 * time.Second))

	expectedVersions := 0
	for n := 0; n < 10; n++ {
		_test.AppendData(file, "more content")
		time.Sleep(time.Duration(1 * time.Second))
		expectedVersions++
	}
	infos, err := ioutil.ReadDir(path.Join(root, _test.Md5(file)))
	_test.FailIfError(err)
	versions := len(infos)
	_test.CleanFolder(root)
	_test.CleanFolder(fileFolder)
	if versions != expectedVersions {
		b.Fatalf("Expected versions are %d , got %d", expectedVersions, versions)
	}
}
