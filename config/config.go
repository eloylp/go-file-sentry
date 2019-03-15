package config

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"
)

type Config struct {
	storagePath string
	wFiles      []string
	socket      *url.URL
}

func (c *Config) Socket() *url.URL {
	return c.socket
}

func (c *Config) WFiles() []string {
	return c.wFiles
}
func (c *Config) StoragePath() string {
	return c.storagePath
}
func NewConfig(storagePath string, watchedFiles []string, socket *url.URL) *Config {
	return &Config{storagePath: storagePath, wFiles: watchedFiles, socket: socket}
}

func NewConfigFromParams() *Config {

	var wFilesPath string
	var storagePath string
	var addr string

	flag.StringVar(&wFilesPath, "f", "", "The path to files list to watch")
	flag.StringVar(&storagePath, "s", "", "The root to the storage to store file versions")
	flag.StringVar(&addr, "l", "unix:///var/run/go-file-sentry.sock", "The socket for listening cli commands")
	flag.Parse()

	filesWatched, err := parseWFiles(wFilesPath)
	if err != nil {
		log.Fatal(err)
	}
	listen, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	return NewConfig(storagePath, filesWatched, listen)
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
