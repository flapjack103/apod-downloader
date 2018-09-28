package main

import (
	"sync"

	"github.com/DataDog/datadog-go/statsd"
)

// Pool defines a worker pool for handling tasks concurrently. Jobs are put
// on the download channel and consumers by worker goroutines.
type Pool struct {
	download chan string
	exit     chan bool
	wg       sync.WaitGroup
	stats    *statsd.Client
}

// NewPool initializes a pool of n workers
func NewPool(n int, stats *statsd.Client) *Pool {
	p := &Pool{
		download: make(chan string, 1024),
		stats:    stats,
	}

	p.wg.Add(n)
	for i := 0; i < n; i++ {
		d := NewDownloader(p.download, stats)
		go func() {
			d.Run()
			p.wg.Done()
		}()
	}

	return p
}

// Download queues the given url for downloading
func (p *Pool) Download(url string) {
	p.download <- url
}

// Close shuts down the pool and waits for all workers to exit
func (p *Pool) Close() {
	close(p.download)
	p.wg.Wait()
}
