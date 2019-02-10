package _test

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/nu7hatch/gouuid"
	"log"
	"os"
	"path/filepath"
)

func WriteFile(folderPath string, name string, content string) string {
	path := filepath.Join(folderPath, string(os.PathSeparator), name)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	FailIfError(err)
	defer file.Close()
	_, err = file.WriteString(content)
	FailIfError(err)
	return path
}

func AppendData(path string, content string) {
	testFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
	FailIfError(err)
	defer testFile.Close()
	_, err = testFile.WriteString(content)
	FailIfError(err)
}

func FailIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const testRootFolderName string = "tmp"
const testFolderPrefix string = "go_file_sentry_test_"

func CreateTestStorageFolder() string {
	uuidGen, _ := uuid.NewV4()
	folder := filepath.Join(string(os.PathSeparator), testRootFolderName, testFolderPrefix+uuidGen.String())
	_ = os.Mkdir(folder, 0755)
	return folder
}

func CreateFixedTestFolder(suffix string) string {
	folder := filepath.Join(string(os.PathSeparator), testRootFolderName, testFolderPrefix+suffix)
	_ = os.Mkdir(folder, 0755)
	return folder
}

func CleanFolder(path string) {
	err := os.RemoveAll(path)
	FailIfError(err)
}

func Md5(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	containerName := hex.EncodeToString(hasher.Sum(nil))
	return containerName
}

func FsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func GetTestResource(name string) string {
	const testResourceDir = "test"
	return filepath.Join(testResourceDir, name)
}
