package main

import (
	"code.google.com/p/goauth2/oauth"
	"github.com/digitalocean/godo"
)

// create a new client using the provided access token
func NewClient(token string) *Client {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}

	return &Client{godo.NewClient(t.Client())}
}

// an invisible wrapper around godo.Client to allow for method additions
type Client struct {
	*godo.Client
}
