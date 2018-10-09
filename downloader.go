package main

import (
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/pkg/errors"

	log "github.com/cihub/seelog"
)

// Downloader defines a type of worker for downloading images given a URL
type Downloader struct {
	jobs  chan string
	stats *statsd.Client
}

// NewDownloader creates a worker for downloading images
func NewDownloader(jobs chan string, stats *statsd.Client) *Downloader {
	return &Downloader{
		jobs:  jobs,
		stats: stats,
	}
}

// Run starts the downloader listening for jobs
func (d *Downloader) Run() {
	for imageURL := range d.jobs {
		//donwload the url
		now := time.Now()
		if err := d.download(imageURL); err != nil {
			log.Errorf("Worker encountered error: %s\n", err)
		}

		// collect some stats
		sec := time.Now().Sub(now).Seconds()
		log.Infof("downloaded %s in %.2fs", imageURL, sec)
		// d.stats.Incr("download.count", nil, 1)
		// d.stats.Gauge("download.duration", sec, nil, 1)
	}
}

func (d *Downloader) download(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "downloading image")
	}

	defer resp.Body.Close()

	//open a file for writing
	file, err := os.Create(path.Join(ImageDir, path.Base(url)))
	if err != nil {
		return errors.Wrap(err, "creating image file")
	}

	// copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.Wrap(err, "copying image data")
	}

	return errors.Wrap(file.Close(), "closing image file")
}
