package main

import (
	"fmt"
	"os"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	log "github.com/cihub/seelog"
)

const (
	// ImageDir sets the download location for our image data
	ImageDir = "/tmp/apod/images"
	// AppNamespace is the prefix for all metrics
	AppNamespace = "apod"
)

// init the statsd client
// more info here: https://github.com/DataDog/datadog-go
func initStatsd() (*statsd.Client, error) {
	c, err := statsd.New("127.0.0.1:8125")
	if err != nil {
		return nil, err
	}
	c.Namespace = AppNamespace
	return c, nil
}

func main() {
	fmt.Println("Getting started")

	// defaults
	downloadCount := 10
	workers := 1
	apiKey := os.Getenv("NASA_API_KEY")

	stats, err := initStatsd()
	if err != nil {
		log.Criticalf("Error initializing stats client: %s", err)
		return
	}

	cli := NewAPODClient(apiKey)
	urls, err := cli.FetchImageURLs(downloadCount)
	if err != nil {
		log.Criticalf("Error fetching image urls: %s", err)
		return
	}

	start := time.Now()
	p := NewPool(workers, stats)
	for _, url := range urls {
		p.Download(url)
	}
	p.Close()

	td := time.Now().Sub(start).Seconds()
	stats.Gauge("fetch.duration", td, nil, 1)
}
