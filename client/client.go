package client

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:4000"

// create a client type that wraps a the http.client
type Client struct {
	*http.Client
}

// create a new client helper mehtod
func New() Client {
	return Client{
		&http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

// Create a Method that will handle all the request to our client.
// This Method will return an error (if any), response code, response header, and response body
func (c *Client) NewRequest(method string, url string, body io.Reader, headers http.Header) (error, int, http.Header, string) {
	req, err := http.NewRequest(method, baseURL+url, body)
	if err != nil {
		return err, 0, http.Header{}, ""
	}

	//add the headers (if provided) to the request
	for k, v := range headers {
		req.Header[k] = v
	}

	resp, err := c.Do(req)
	if err != nil {
		return err, 0, http.Header{}, ""
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, 0, http.Header{}, ""
	}
	respBody = bytes.TrimSpace(respBody)

	return nil, resp.StatusCode, resp.Header, string(respBody)
}
