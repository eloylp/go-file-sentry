package diff_test

import (
	"github.com/eloylp/go-file-sentry/diff"
	"github.com/eloylp/go-file-sentry/file"
	"path/filepath"
	"testing"
)

func getTestResource(resourceName string) string {
	return filepath.Join(resourceName)
}

func TestGetDiffOfFiles(t *testing.T) {

	fileA := file.File{}
	fileA.Path = getTestResource("fileA.conf")
	fileB := file.File{}
	fileB.Path = getTestResource("fileB.conf")
	expectedDiffFile := file.File{}
	expectedDiffFile.Path = getTestResource("expected.diff")

	diffOfFiles := diff.GetDiffOfFiles(fileA, fileB)
	expectedDiffOfFiles := string(expectedDiffFile.GetData())

	if expectedDiffOfFiles != diffOfFiles {
		t.Errorf("Files diff do not match expected result.")
	}
}