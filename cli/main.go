package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
)

func main() {

	param := os.Getenv("API_URL")
	if param == "" {
		log.Fatal("API_URL needed.")
	}
	addr, err := url.Parse(param)
	if err != nil {
		log.Fatal(err)
	}
	var c *http.Client
	var prefix string
	if addr.Scheme == "unix" {
		c = unixClient(addr.Path)
		prefix = "http://unix"
	} else {
		c = tcpClient()
		prefix = addr.String()
	}
	r, err := c.Get(prefix + "/status")
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(bytes))
}

func tcpClient() *http.Client {
	return &http.Client{}
}

func unixClient(path string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", path)
			},
		},
	}
}
