package config

import (
	"bufio"
	"flag"
	"log"
	"os"
)

type Config struct {
	storagePath string
	wFiles      []string
}

func (c *Config) WFiles() []string {
	return c.wFiles
}
func (c *Config) StoragePath() string {
	return c.storagePath
}
func NewConfig(storagePath string, watchedFiles []string) *Config {
	return &Config{storagePath: storagePath, wFiles: watchedFiles}
}

func NewConfigFromParams() *Config {
	var wFilesPath string
	var storagePath string
	flag.StringVar(&wFilesPath, "files", "", "The path to files list to watch")
	flag.StringVar(&storagePath, "storage", "", "The root to the storage to store file versions")
	flag.Parse()
	filesWatched, err := parseWFiles(wFilesPath)
	if err != nil {
		log.Fatal(err)
	}
	return NewConfig(storagePath, filesWatched)
}

func parseWFiles(listPath string) ([]string, error) {

	file, err := os.Open(listPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var wFiles []string
	for scanner.Scan() {
		wFiles = append(wFiles, scanner.Text())
	}
	return wFiles, nil
}
