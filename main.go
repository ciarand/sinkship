package main

import (
	"errors"
	"flag"
	"os"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"

	"code.google.com/p/goauth2/oauth"
	"github.com/digitalocean/godo"
)
import "io/ioutil"

func main() {
	// try a couple of different places to find the file
	pat, err := getTokenChain(getTokenFromCli, getTokenFromEnv, getTokenFromFile)
	if err != nil {
		log.Fatalf("couldn't get token: %s", err.Error())
	}

	client := NewClient(pat)

	droplets, _, err := client.Droplets.List(nil)
	if err != nil {
		log.Fatalf("couldn't get a list of your droplets: %s", err.Error())
	}

	if len(droplets) == 0 {
		log.Info("You've got no droplets in your account.")
		return
	}

	log.Infof("Found %d droplets, preparing to delete", len(droplets))
	wg := &sync.WaitGroup{}
	mut := &sync.Mutex{}
	count := 0

	for _, v := range droplets {
		wg.Add(1)
		go func(droplet godo.Droplet) {
			_, err := client.Droplets.Delete(droplet.ID)
			mut.Lock()
			count++
			if err != nil {
				log.Errorf("[%d] couldn't delete droplet %s (%d):\n\t%s", count, droplet.Name, droplet.ID, err)
			} else {
				log.Infof("[%d] deleted %s (%d)", count, droplet.Name, droplet.ID)
			}
			mut.Unlock()

			wg.Done()
		}(v)
	}

	wg.Wait()
}

func NewClient(token string) *Client {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}

	return &Client{godo.NewClient(t.Client())}
}

type Client struct {
	*godo.Client
}

type TokenGetter func() (string, error)

func getTokenChain(getters ...TokenGetter) (string, error) {
	errs := make([]string, len(getters))

	for _, g := range getters {
		str, err := g()
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}

		if str != "" {
			return str, nil
		}
	}

	return "", errors.New(strings.Join(errs, "\n"))
}

func getTokenFromFile() (string, error) {
	bytes, err := ioutil.ReadFile("token")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func getTokenFromEnv() (string, error) {
	return os.Getenv("DO_TOKEN"), nil
}

func getTokenFromCli() (string, error) {
	var str *string

	flag.StringVar(str, "token", "", "The token to use with the DO API")
	flag.Parse()

	return *str, nil
}
