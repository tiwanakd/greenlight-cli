package client

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://greenlight-api.com"

// create a client type that wraps a the http.client
type Client struct {
	*http.Client
}

// create a response struch that will be returned by the Client
type Response struct {
	Code int
	Body string
	Err  error
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
// This Method will return the Response struch
func (c *Client) NewRequest(method string, url string, body io.Reader, headers http.Header) Response {
	req, err := http.NewRequest(method, baseURL+url, body)
	if err != nil {
		return Response{http.StatusBadRequest, "", err}
	}

	//add the headers (if provided) to the request
	for k, v := range headers {
		req.Header[k] = v
	}

	resp, err := c.Do(req)
	if err != nil {
		return Response{http.StatusBadRequest, "", err}
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{http.StatusBadRequest, "", err}
	}
	respBody = bytes.TrimSpace(respBody)

	return Response{resp.StatusCode, string(respBody), nil}
}
