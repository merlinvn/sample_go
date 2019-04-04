package net

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL *url.URL
	client  *http.Client
}

type RequestHeader struct {
	ContentType string
	Accept      string
	UserAgent   string
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

func (client *Client) NewRequest(method string, header *RequestHeader, urlString string, body interface{}) (*http.Request, error) {
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

	setHeader(req, header, body)

	return req, nil

}

func deferClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

func (client *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := ctxhttp.Do(ctx, client.client, req)

	if err != nil {
		return nil, err
	}

	defer deferClose(resp.Body)

	err = CheckResponse(resp)

	if err != nil {
		return nil, err
	}

	switch v.(type) {
	case *interface{}:
		data, err := ioutil.ReadAll(resp.Body)
		if err == nil && data != nil {
			err := json.Unmarshal(data, v)
			if err != nil {
				return nil, err
			}
		}
	default:
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}

type ErrorResponse struct {
	*http.Response
	Message string `json:"message"`
}

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v", er.Response.Request.Method, er.Response.Request.URL.Path, er.StatusCode, er.Message)
}

func CheckResponse(response *http.Response) error {
	if response.StatusCode < 300 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: response}
	data, err := ioutil.ReadAll(response.Body)
	if err == nil && data == nil {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			return err
		}
	}

	return errorResponse
}

func setHeader(req *http.Request, header *RequestHeader, body interface{}) {
	if header == nil {
		return
	}

	if body != nil {
		checkAndSetHeaderValue(req, "Content-Type", header.ContentType)
	}
	checkAndSetHeaderValue(req, "Accept", header.Accept)
	checkAndSetHeaderValue(req, "User-Agent", header.UserAgent)
}

func checkAndSetHeaderValue(req *http.Request, headerKey string, headerValue string, ) {
	if headerValue != "" {
		req.Header.Set(headerKey, headerValue)
	}
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
