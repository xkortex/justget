package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Config struct {
	timeout    time.Duration
	targetUrls []string
}

func parseArgs() (cfg Config) {
	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), `NAME
	justget
SYNOPSIS
	justget [options...] URL [URLs...]

DESCRIPTION
	Issue GET requests to one or more http URLs
`)
		flag.PrintDefaults()
	}

	flag.DurationVar(&cfg.timeout, "timeout", time.Duration(1e9), "Timeout duration string")
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatalln("No URLs specified")
	}
	cfg.targetUrls = flag.Args()
	return
}

func get(url string, c chan string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	c <- string(body)
}

func main() {
	cfg := parseArgs()
	c := make(chan string)
	for _, url := range cfg.targetUrls {
		go get(url, c)
	}
	select {
	case res := <-c:
		fmt.Println(res)
	case <-time.After(cfg.timeout):
		log.Fatalf("Timed out after %v\n", cfg.timeout)
	}

}
