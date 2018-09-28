package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// Client for interacting with NASA's APOD (Astronomy Picture of the Day) API

const (
	// BaseURL is the APOD API URL
	BaseURL = "https://api.nasa.gov/planetary/apod"
	// APODTypeImage for defining image media type
	APODTypeImage = "image"
	// APODTypeVideo for defining video media type
	APODTypeVideo = "video"
	// APODDateFormat defines the date formating for an API request
	APODDateFormat = "2006-01-02"
)

// APODImageMeta defines the response returned by the APOD API for a single
// image object. This only defines a subset of the returned fields that we
// care about. Full information about the API and available fields can be
// found at https://github.com/nasa/apod-api
type APODImageMeta struct {
	URL       string `json:"hdurl"`
	MediaType string `json:"media_type"`
	Date      string `json:"date"`
}

// APODClient defines the client for interacting with the API
type APODClient struct {
	URL    string
	APIKey string
}

// NewAPODClient creates a new client for communicating with APOD API
func NewAPODClient(apiKey string) *APODClient {
	return &APODClient{
		URL:    BaseURL,
		APIKey: apiKey,
	}
}

// FetchImageURLs returns a list of all image URLs from nDaysAgo
func (c *APODClient) FetchImageURLs(nDaysAgo int) ([]string, error) {
	startDate := time.Now().AddDate(0, 0, -nDaysAgo)

	// make the request
	resp, err := http.Get(c.buildURL(startDate))
	if err != nil {
		return nil, errors.Wrap(err, "error fetching data from APOD API")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Received non-200 status code %d", resp.StatusCode)
	}

	// parse the response
	var imageMeta []APODImageMeta
	if err := json.NewDecoder(resp.Body).Decode(&imageMeta); err != nil {
		return nil, errors.Wrap(err, "error parsing API response")
	}

	var urls []string
	for _, img := range imageMeta {
		if img.MediaType != APODTypeImage {
			continue
		}
		urls = append(urls, img.URL)
	}

	return urls, nil
}

// helper method for constructing the API url for the given start date
func (c *APODClient) buildURL(startDate time.Time) string {
	v := url.Values{}
	v.Add("api_key", c.APIKey)
	v.Add("start_date", startDate.Format(APODDateFormat))
	return fmt.Sprintf("%s?%s", c.URL, v.Encode())
}
