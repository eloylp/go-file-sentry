package scan

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/eloylp/go-file-sentry/file"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ScanFile(path string) file.File {
	return FileInfoGatherer(path)
}

func FQDNCalculator(file *file.File) {

	const pathSeparator string = "_"
	const sysDirNameSeparator = "-"
	const systemDirNameDatePart string = "20060102150405"

	filePathPart := strings.Replace(file.Path, string(filepath.Separator), pathSeparator, -1)
	if strings.HasPrefix(filePathPart, pathSeparator) {
		filePathPart = strings.TrimLeft(filePathPart, pathSeparator)
	}
	fileDatePart := file.Time.Format(systemDirNameDatePart)
	parts := []string{filePathPart, file.Sum, fileDatePart}
	file.FQDN = strings.Join(parts, sysDirNameSeparator)
}

func FileInfoGatherer(path string) file.File {
	targetFile := file.File{}
	readFile, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	targetFile.Time = readFile.ModTime()
	targetFile.Path = path
	targetFile.Sum = GetFileSum(path)
	FQDNCalculator(&targetFile)
	return targetFile
}

func GetFileSum(path string) string {
	targetFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer targetFile.Close()
	sum := md5.New()
	if _, err = io.Copy(sum, targetFile); err != nil {
		log.Fatal(err)
	}
	sumInBytes := sum.Sum(nil)[:16]
	sumString := hex.EncodeToString(sumInBytes)
	return sumString
}
