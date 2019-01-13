package diff_test

import (
	"github.com/eloylp/go-file-sentry/_test"
	"github.com/eloylp/go-file-sentry/diff"
	"github.com/eloylp/go-file-sentry/file"
	"testing"
)

func TestGetDiffOfFiles(t *testing.T) {

	fileA := file.NewFile(_test.GetTestResource("fileA.conf"))
	fileB := file.NewFile(_test.GetTestResource("fileB.conf"))
	expectedDiffFile := file.NewFile(_test.GetTestResource("expected.diff"))

	diffOfFiles := diff.GetDiffOfFiles(fileA, fileB)
	expectedDiffOfFiles := string(expectedDiffFile.Data())

	if expectedDiffOfFiles != diffOfFiles {
		t.Errorf("Files diff do not match expected result.")
	}
}
