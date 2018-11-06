package scan

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ScanFile(path string) File {
	return FileInfoGatherer(path)
}

func FQDNCalculator(file *File) {

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

func FileInfoGatherer(path string) File {
	file := File{}
	readFile, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	file.Time = readFile.ModTime()
	file.Path = path
	file.Sum = GetFileSum(path)
	FQDNCalculator(&file)
	return file
}

func GetFileSum(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	sum := md5.New()
	if _, err = io.Copy(sum, file); err != nil {
		log.Fatal(err)
	}
	sumInBytes := sum.Sum(nil)[:16]
	sumString := hex.EncodeToString(sumInBytes)
	return sumString
}
