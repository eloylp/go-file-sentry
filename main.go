package main

import (
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/sentry"
)

func main() {

	done := make(chan bool)
	cfg := config.NewConfig(
		"/tmp/test/root",
		[]string{"/tmp/test/file1.txt", "/tmp/test/file1.txt"})

	sentry.StartSentry(cfg)
	<-done
}
