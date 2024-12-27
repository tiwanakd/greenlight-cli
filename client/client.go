package client

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:4000"

type Client struct {
	*http.Client
}

func New() Client {
	return Client{
		&http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

func (c *Client) GetRequest(url string) (error, int, http.Header, string) {
	resp, err := c.Get(baseURL + url)
	if err != nil {
		return err, 0, http.Header{}, ""
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, 0, http.Header{}, ""
	}
	body = bytes.TrimSpace(body)

	return nil, resp.StatusCode, resp.Header, string(body)
}

func (c *Client) PostRequest(url string, body io.Reader) (error, int, http.Header, string) {
	resp, err := c.Post(baseURL+url, "application/json", body)
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
