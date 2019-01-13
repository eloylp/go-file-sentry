package file_test

import (
	"github.com/eloylp/go-file-sentry/_test"
	"github.com/eloylp/go-file-sentry/file"
	"testing"
)

func TestFileGetData(t *testing.T) {
	f := file.NewFile(_test.GetTestResource("file.txt"))
	data := f.Data()
	text := string(data)
	expectedContent := "This is a test file."
	if text != expectedContent {
		t.Errorf("Expected content of file is %s , result was %s", expectedContent, text)
	}
}

func TestFileGetName(t *testing.T) {
	f := file.NewFile(_test.GetTestResource("file.txt"))
	name := f.Name()
	expectedName := "file.txt"
	if name != expectedName {
		t.Errorf("Expected name of file is %s , result was %s", expectedName, name)
	}
}

func TestNewFile(t *testing.T) {
	f := file.NewFile(_test.GetTestResource("file.txt"))
	expectedSum := "3de8f8b0dc94b8c2230fab9ec0ba0506"
	if f.Sum() != expectedSum {
		t.Errorf("Expected sum of file is %s , result was %s", expectedSum, f.Sum())
	}
	expectedFQDN := "3de8f8b0dc94b8c2230fab9ec0ba0506-20181112183838"
	if f.FQDN() != expectedFQDN {
		t.Errorf("Expected fqdn of file is %s , result was %s", expectedSum, f.FQDN())
	}

}
