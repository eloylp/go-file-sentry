package _test

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/nu7hatch/gouuid"
	"log"
	"os"
	"path/filepath"
)

func WriteFile(testFolderPath string, name string, content string) string {
	filePath := filepath.Join(testFolderPath, string(os.PathSeparator), name)
	testFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	FailIfError(err)
	defer testFile.Close()
	_, err = testFile.WriteString(content)
	FailIfError(err)
	return filePath
}

func AppendData(filePath string, content string) {
	testFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
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
	testFolder := filepath.Join(string(os.PathSeparator), testRootFolderName, testFolderPrefix+uuidGen.String())
	_ = os.Mkdir(testFolder, 0755)
	return testFolder
}

func CreateFixedTestFolder(suffix string) string {
	testFolder := filepath.Join(string(os.PathSeparator), testRootFolderName, testFolderPrefix+suffix)
	_ = os.Mkdir(testFolder, 0755)
	return testFolder
}

func CleanFolder(path string) {
	err := os.RemoveAll(path)
	FailIfError(err)
}

func CalculateMd5(input string) string {
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

func GetTestResource(resourceName string) string {
	const testResourceDir = "test"
	return filepath.Join(testResourceDir, resourceName)
}
