package main

import (
	"context"
	"fmt"
	"github.com/merlinvn/sample_go/net"
	"net/http"
	"os"
)

var apiToken = os.Getenv("GITHUB_API_TOKEN")

// Repository represents a GitHub repository.
type Repository struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GitURL      string `json:"git_url"`
}

type GraphQLBody struct{
	query string
}

func main() {
	ctx := context.Background()
	c := net.NewClient(ctx, apiToken, "https://api.github.com/graphql")

	header := &net.RequestHeader{ContentType: "application/json", Accept: "application/json", UserAgent: "My UserAgent"}

	body := map[string]string{"query": "{viewer{login name}}"}
	//`{"query":"{viewer{login name}}"}`

	req, err := c.NewRequest(http.MethodPost, header, "", body)

	if err != nil {
		fmt.Println(err)
	}

	//var result []*Repository
	var result interface{}

	_, err = c.Do(ctx, req, &result)

	if err != nil {
		fmt.Println(err)
	}

	m := result.(map[string]interface{})
	for key, value := range m {
		fmt.Println("Key:", key, "Value:", value)
	}
}

func github_v3() {
	ctx := context.Background()
	c := net.NewClient(ctx, apiToken, "https://api.github.com/")

	header := &net.RequestHeader{ContentType: "application/json", Accept: "application/vnd.github.v3+json", UserAgent: "My UserAgent"}

	req, err := c.NewRequest("GET", header, "user/repos", "")

	if err != nil {
		fmt.Println(err)
	}

	//var result []*Repository
	var result interface{}

	_, err = c.Do(ctx, req, &result)

	if err != nil {
		fmt.Println(err)
	}

	m := result.([]interface{})
	fmt.Println(len(m))
	//for k, v := range m {
	//	switch vv := v.(type) {
	//	case string:
	//		fmt.Println(k, "is string", vv)
	//	case float64:
	//		fmt.Println(k, "is float64", vv)
	//	case []interface{}:
	//		fmt.Println(k, "is an array:")
	//		for i, u := range vv {
	//			fmt.Println(i, u)
	//		}
	//	default:
	//		fmt.Println(k, "is of a type I don't know how to handle")
	//	}
	//}
	//fmt.Println(len(result.(map[string]interface{})))
}
