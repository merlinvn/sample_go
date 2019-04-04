package main

import (
	"context"
	"github.com/merlinvn/sample_go/net"
	"os"
)

var apiToken = os.Getenv("GITHUB_API_TOKEN")

func main(){
	ctx := context.Background()
	c := net.NewClient(ctx, apiToken, "https://api.github.com/graphql")


}
