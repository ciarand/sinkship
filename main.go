package main

import (
	"sync"

	log "github.com/Sirupsen/logrus"

	"code.google.com/p/goauth2/oauth"
	"github.com/digitalocean/godo"
)
import "io/ioutil"

func main() {
	pat, err := getTokenFromFile("token")
	if err != nil {
		log.Fatalf("couldn't get token: %s", err.Error())
	}

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: pat},
	}

	client := godo.NewClient(t.Client())

	drops, _, err := client.Droplets.List(nil)
	if err != nil {
		log.Fatal("couldn't get a list of your droplets: %s", err.Error())
	}

	if len(drops) == 0 {
		log.Info("You've got no droplets in your account.")
		return
	}

	log.Info("Found %d droplets, preparing to delete", len(drops))
	wg := &sync.WaitGroup{}

	for _, d := range drops {
		wg.Add(1)
		go func() {
			err := tryToDeleteDroplet(client, d)
			if err != nil {
				log.Errorf("couldn't delete droplet %s (%d): %s", d.Name, d.ID, err)
			}

			log.Info("deleted %s (%s)", d.Name, d.ID)
			wg.Done()
		}()
	}

	wg.Wait()
}

func getTokenFromFile(file string) (string, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func tryToDeleteDroplet(c *godo.Client, d godo.Droplet) error {
	_, err := c.Droplets.Delete(d.ID)

	if err != nil {
		return err
	}

	return nil
}
