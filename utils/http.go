package utils

import (
	"io"
	"log"
	"net/http"
	"time"
)

func HttpRequest(method, url string, headers map[string]string, body io.Reader) *http.Response {
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest(method, url, body)

	if err != nil {
		log.Print(err)
		return nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:124.0) Gecko/20100101 Firefox/124.0")

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	// Try 5 times then bail out
	for i := 1; i <= 5; i++ {
		response := makeRequest(client, req)

		if response != nil {
			return response
		}

		time.Sleep(time.Second * 5 * time.Duration(i))
	}

	return nil
}

func makeRequest(client *http.Client, req *http.Request) *http.Response {
	response, err := client.Do(req) //#nosec G704 -- URLs are from user config

	if err != nil {
		return nil
	}

	return response
}
