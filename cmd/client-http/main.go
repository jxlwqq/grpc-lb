package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"io/ioutil"
	"net/http"
	"time"
)

var url = flag.String("url", "", "server url")

const (
	// disable Keep-Alive to act as massive independent clients
	disableKeepalive = true
	timeout          = 5 * time.Second
	waitDuration     = 1 * time.Second
)

var counter int

func connectAndShowResponse(client *http.Client, url string) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	counter++
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("#%d: %s\n", counter, string(body))
}

func main() {
	flag.Parse()

	client := &http.Client{
		Transport: &http.Transport{
			// disable Keep-Alive
			// @see https://www.cnblogs.com/cobbliu/p/4517598.html
			// @see https://nanxiao.me/en/a-brief-intro-of-tcp-keep-alive-in-gos-http-implementation/
			DisableKeepAlives: disableKeepalive,
		},
		Timeout: timeout,
	}

	for {
		connectAndShowResponse(client, *url)
		time.Sleep(waitDuration)
	}
}
