package diff

import (
	. "github.com/eloylp/go-file-sentry/file"
	. "github.com/pmezard/go-difflib/difflib"
	"log"
)

func GetDiffOfFiles(file File, file2 File) string {

	diff := ContextDiff{
		A:        SplitLines(string(file.GetData())),
		B:        SplitLines(string(file2.GetData())),
		FromFile: "Original",
		ToFile:   "Current",
		Context:  3,
		Eol:      "\n",
	}

	result, err := GetContextDiffString(diff)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
