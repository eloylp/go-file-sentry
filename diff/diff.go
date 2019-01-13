package diff

import (
	"github.com/eloylp/go-file-sentry/file"
	"github.com/pmezard/go-difflib/difflib"
	"log"
)

func GetDiffOfFiles(file *file.File, file2 *file.File) string {

	diff := difflib.ContextDiff{
		A:        difflib.SplitLines(string(file.Data())),
		B:        difflib.SplitLines(string(file2.Data())),
		FromFile: "Original",
		ToFile:   "Current",
		Context:  3,
		Eol:      "\n",
	}

	result, err := difflib.GetContextDiffString(diff)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
