package github

import (
	"context"
	"fmt"
	"net/http"
)

const (
	reposPathWithUser        = "users/%v/repos"
	graphqlReposPathWithUser = ``
	defaultReposPath = "user/repos"
)

// RepositoriesService handles all the repositories actions

// Github API docs: https://docs.gitlab.com/ce/api/repositories.html
type RepositoriesService struct {
	client *Client
}

// Repository represents a GitHub repository.
type Repository struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GitURL      string `json:"git_url"`
}

// List the repositories for a user.
//
// GitHub API docs: https://developer.github.com/v3/repos/#list-user-repositories
func (s *RepositoriesService) List(ctx context.Context, user string) ([]*Repository, *http.Response, error) {
	var path string
	if user != "" {
		path = fmt.Sprintf(reposPathWithUser, user)
	} else {
		path = defaultReposPath
	}

	//body:= `
	//{
	//	query: query { viewer { login }}
	//}
	//`

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var repos []*Repository
	resp, err := s.client.Do(ctx, req, &repos)

	if err != nil {
		return nil, resp, err
	}

	return repos, resp, nil
}


//?query={
//user(login: "%v") {
//repositories(first: 30) {
//totalCount
//pageInfo {
//endCursor
//hasNextPage
//}
//nodes {
//id,
//name,
//description,
//sshUrl
//}
//}
//}
//}