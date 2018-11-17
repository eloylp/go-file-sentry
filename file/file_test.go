package file_test

import (
	"github.com/eloylp/go-file-sentry/file"
	"path/filepath"
	"testing"
)

func TestFileGetData(t *testing.T) {
	f := file.File{}
	f.Path = filepath.Join("file.txt")
	data := f.GetData()
	text := string(data)
	expectedContent := "This is a test file."
	if text != expectedContent {
		t.Errorf("Expected content of file is %s , result was %s", expectedContent, text)
	}
}
