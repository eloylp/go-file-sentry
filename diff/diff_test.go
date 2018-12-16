package diff_test

import (
	"github.com/eloylp/go-file-sentry/diff"
	"github.com/eloylp/go-file-sentry/file"
	"testing"
)

func TestGetDiffOfFiles(t *testing.T) {

	fileA := file.NewFile("fileA.conf")
	fileB := file.NewFile("fileB.conf")
	expectedDiffFile := file.NewFile("expected.diff")

	diffOfFiles := diff.GetDiffOfFiles(fileA, fileB)
	expectedDiffOfFiles := string(expectedDiffFile.GetData())

	if expectedDiffOfFiles != diffOfFiles {
		t.Errorf("Files diff do not match expected result.")
	}
}
