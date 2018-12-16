package file_test

import (
	"github.com/eloylp/go-file-sentry/file"
	"path/filepath"
	"testing"
)

func TestFileGetData(t *testing.T) {
	f := file.NewFile(filepath.Join("file.txt"))
	data := f.GetData()
	text := string(data)
	expectedContent := "This is a test file."
	if text != expectedContent {
		t.Errorf("Expected content of file is %s , result was %s", expectedContent, text)
	}
}

func TestFileGetName(t *testing.T) {
	f := file.NewFile(filepath.Join("file.txt"))
	name := f.GetName()
	expectedName := "file.txt"
	if name != expectedName {
		t.Errorf("Expected name of file is %s , result was %s", expectedName, name)
	}
}

func TestNewFile(t *testing.T) {
	f := file.NewFile(filepath.Join("file.txt"))
	expectedSum := "3de8f8b0dc94b8c2230fab9ec0ba0506"
	if f.GetSum() != expectedSum {
		t.Errorf("Expected sum of file is %s , result was %s", expectedSum, f.GetSum())
	}
	expectedFQDN := "3de8f8b0dc94b8c2230fab9ec0ba0506-20181112183838"
	if f.GetFQDN() != expectedFQDN {
		t.Errorf("Expected FQDN of file is %s , result was %s", expectedSum, f.GetFQDN())
	}

}
