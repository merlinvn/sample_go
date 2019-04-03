package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	defaultBaseUrl      = "https://api.github.com/graphql"
	acceptVersionHeader = "Accept: application/vnd.github.v4.idl"
)

type Client struct {
	client  *http.Client
	baseUrl *url.URL
	// Services used for handling the Github API.
	Repositories *RepositoriesService
}

func (client *Client) SetBaseUrl(urlString string) error {
	baseUrl, err := url.Parse(urlString)
	if err != nil {
		return err
	}

	client.baseUrl = baseUrl
	return nil
}

func NewClient(ctx context.Context, token string) *Client {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := newClient(tc)
	return client
}

func newClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{client: httpClient}
	c.SetBaseUrl(defaultBaseUrl)

	// Create all the public services.
	c.Repositories = &RepositoriesService{client:c}

	return c

}

func (client *Client) NewRequest(method, urlString string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlString)

	if err != nil {
		return nil, err
	}

	u := client.baseUrl.ResolveReference(rel)
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

func (client *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := ctxhttp.Do(ctx, client.client, req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = CheckResponse(resp)

	if err != nil {
		return nil, err
	}

	err = client.decodeResponse(resp.Body, v)

	if err != nil {
		return nil, err
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
		json.Unmarshal(data, errorResponse)
	}

	return errorResponse
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

func (client *Client) decodeResponse(body io.ReadCloser, v interface{}) error {
	if v != nil {
		var err error
		err = json.NewDecoder(body).Decode(v)

		return err
	}

	return nil
}
