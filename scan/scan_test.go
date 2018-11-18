package scan_test

import (
	"github.com/eloylp/go-file-sentry/file"
	"github.com/eloylp/go-file-sentry/scan"
	"path/filepath"
	"testing"
	"time"
)

func getTestResource(resourceName string) string {
	return filepath.Join(resourceName)
}

func TestFQDNCalculator(t *testing.T) {

	fileTime, err := time.Parse("2006-01-02 15:04:05", "2018-01-01 13:43:54")
	if err != nil {
		t.Errorf(err.Error())
	}
	testFile := file.File{}
	testFile.Path = "/etc/fstab"
	testFile.Sum = "587399e23181c0a8862b1c8c2a2225a6"
	testFile.Time = fileTime
	scan.FQDNCalculator(&testFile)
	expectedFQDN := "etc_fstab-587399e23181c0a8862b1c8c2a2225a6-20180101134354"

	if testFile.FQDN != expectedFQDN {
		t.Errorf("Invalid FQDN. Was detected \"%s\" and expected is \"%s\"", testFile.FQDN, expectedFQDN)
	}
}

func TestFileInfoGatherer(t *testing.T) {
	rawFilePath := getTestResource("fstab")
	testFile := scan.FileInfoGatherer(rawFilePath)
	expectedFQDN := "fstab-587399e23181c0a8862b1c8c2a2225a6-20181029183845"
	expectedTime := "2018-10-29 18:38:45"
	expectedSum := "587399e23181c0a8862b1c8c2a2225a6"

	if testFile.Path != rawFilePath {
		t.Errorf("Invalid path. Was detected \"%s\" and expected is \"%s\"", testFile.Path, rawFilePath)
	}
	if testFile.FQDN != expectedFQDN {
		t.Errorf("Invalid FQDN. Was detected \"%s\" and expected is \"%s\"", testFile.FQDN, expectedFQDN)
	}

	if testFile.Sum != expectedSum {
		t.Errorf("Invalid testFile sum. Was detected \"%s\" and expected is \"%s\"", testFile.Sum, expectedSum)
	}

	fileTime := testFile.Time.Format("2006-01-02 15:04:05")
	if fileTime != expectedTime {
		t.Errorf("Invalid testFile mod Time. Was detected \"%s\" and expected is \"%s\"", fileTime, expectedTime)
	}
}
