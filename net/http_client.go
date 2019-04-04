package net

import (
	"bytes"
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL   *url.URL
	userAgent string

	client *http.Client
}

func (client *Client) SetBaseUrl(urlString string) error {
	baseUrl, err := url.Parse(urlString)
	if err != nil {
		return err
	}
	client.baseURL = baseUrl
	return nil
}

func NewClient(ctx context.Context, token string, baseURL string) *Client {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := newClient(tc, baseURL)
	return client
}

func newClient(httpClient *http.Client, baseURL string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{client: httpClient}

	err := c.SetBaseUrl(baseURL)
	if err != nil {
		return nil
	}
	return c
}

func (client *Client) NewRequest(method, urlString string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlString)

	if err != nil {
		return nil, err
	}
	u := client.baseURL.ResolveReference(rel)

	buff, err := client.encodeBody(body)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), buff)

	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", acceptVersionHeader)

	return req, nil

}

func (client *Client) encodeBody(body interface{}) (io.ReadWriter, error) {
	var buff io.ReadWriter

	if body != nil {
		buff = new(bytes.Buffer)
		enc := json.NewEncoder(buff)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)

		if err != nil {
			return nil, err
		} else {
			return buff, nil
		}
	}

	return nil, nil
}
