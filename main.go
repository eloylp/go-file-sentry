package main

import (
	"github.com/eloylp/go-file-sentry/config"
	"github.com/eloylp/go-file-sentry/sentry"
)

func main() {

	done := make(chan bool)
	cfg := config.NewConfigFromParams()
	sentry.Start(cfg)
	<-done
}
