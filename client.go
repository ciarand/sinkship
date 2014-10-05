package main

import (
	"code.google.com/p/goauth2/oauth"
	"github.com/digitalocean/godo"
)

// newClient creates a new Client using the provided access token
func newClient(token string) *client {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}

	return &client{godo.NewClient(t.Client())}
}

// client is an invisible wrapper around godo.Client to allow for method
// additions
type client struct {
	*godo.Client
}
