package term

import (
	"log"
	"os"
	"os/signal"
)

func Listen(shutdown chan bool) {
	sg := make(chan os.Signal, 1)
	signal.Notify(sg, os.Interrupt)
	go func() {
		<-sg
		log.Println("Gracefully ending watching ...")
		shutdown <- true
	}()
}
