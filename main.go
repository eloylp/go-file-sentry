package main

import (
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/sentry"
	"github.com/eloylp/go-file-sentry/term"
	"log"
)

func main() {
	shutdown := make(chan bool)
	term.Listen(shutdown)
	cfg := config.NewConfigFromParams()
	sentry.Start(cfg)
	log.Println("Starting watching of files changes ...")
	<-shutdown
}
