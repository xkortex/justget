package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Config struct {
	fail       bool
	silent     bool
	timeout    time.Duration
	targetUrls []string
	cleanUrls  []string
	verbose    bool
}

type ResponseJ struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Url        string `json:"url"`
}

func castResponse(resp *http.Response, urlStr string) ResponseJ {
	return ResponseJ{Status: resp.Status, StatusCode: resp.StatusCode, Url: urlStr}
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
	flag.BoolVar(&cfg.verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&cfg.verbose, "v", false, "")
	flag.BoolVar(&cfg.fail, "fail", false, "Fail  silently  (no  output at all) on server errors. This is mostly done to better enablescripts etc to better deal with failed attempts")
	flag.BoolVar(&cfg.fail, "f", false, "")
	flag.BoolVar(&cfg.silent, "silent", false, "Silent or quiet mode. Don't show progress meter or error messages.")
	flag.BoolVar(&cfg.silent, "s", false, "")
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatalln("No URLs specified")
	}
	cfg.targetUrls = flag.Args()
	standardizeUrls(&cfg)
	return
}

func standardizeUrl(urlStr string, cfg *Config) string {
	urlp, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.verbose {
		fmt.Fprintf(os.Stderr, "%s -> %+v\n", urlStr, *urlp)
	}
	// deal with dumb bug in the basic parser that can't deal with `example.com:80`
	if len(urlp.Opaque) > 0 && len(urlp.Scheme) > 0 && len(urlp.Host) == 0 {
		return standardizeUrl("http://"+urlp.String(), cfg)
	}
	if len(urlp.Scheme) == 0 {
		return standardizeUrl("http://"+urlp.String(), cfg)
	}
	return urlp.String()
}

func standardizeUrls(cfg *Config) {
	cfg.cleanUrls = make([]string, len(cfg.targetUrls))
	for i, urlStr := range cfg.targetUrls {
		cfg.cleanUrls[i] = standardizeUrl(urlStr, cfg)
	}
}

func get(urlStr string, c chan string, cfg *Config) {
	urlp, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}
	if len(urlp.Scheme) == 0 {
		log.Fatal("Unable to infer a correct URL: %s", urlp.String())
	}
	if cfg.verbose {
		fmt.Fprintf(os.Stderr, "get: %+v\n", urlp.String())
	}
	resp, err := http.Get(urlp.String())
	if err != nil {
		log.Fatalln(err)
	}
	if cfg.verbose {
		fmt.Fprintf(os.Stderr, "%+v\n", *resp)
	}
	if cfg.fail {
		if resp.StatusCode != http.StatusOK {
			j, err := json.MarshalIndent(castResponse(resp, urlStr), "", "  ")
			if err != nil {
				panic(err)
			}
			os.Stdout.Write(j)
			os.Stdout.Write([]byte("\n"))
			os.Exit(1)
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	c <- string(body)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := parseArgs()
	c := make(chan string)
	for _, urlStr := range cfg.cleanUrls {
		go get(urlStr, c, &cfg)
	}
	select {
	case res := <-c:
		fmt.Println(res)
	case <-time.After(cfg.timeout):
		log.Fatalf("Timed out after %v\n", cfg.timeout)
	}

}
