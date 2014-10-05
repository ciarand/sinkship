package main

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/digitalocean/godo"
)

// connects to your DO account and deletes all your Droplets
func main() {
	// try a couple of different places to find the file
	pat, err := getTokenFromChain(getTokenFromCli, getTokenFromEnv, getTokenFromFile)
	if err != nil {
		log.Fatalf("couldn't get token: %s", err.Error())
	}

	client := NewClient(pat)

	// get a list of all the droplets
	droplets, _, err := client.Droplets.List(nil)
	if err != nil {
		log.Fatalf("couldn't get a list of your droplets: %s", err.Error())
	}

	if len(droplets) == 0 {
		log.Info("You've got no droplets in your account.")
		return
	}

	log.Infof("Found %d droplets, preparing to delete", len(droplets))
	// the wait group is so we know how many concurrent requests to wait for
	wg := &sync.WaitGroup{}
	// the mutex is so that they print out in order
	mut := &sync.Mutex{}
	// this is the current index value
	count := 0

	for _, v := range droplets {
		wg.Add(1)
		go func(droplet godo.Droplet) {
			_, err := client.Droplets.Delete(droplet.ID)

			// lock the mutex, increment the count
			mut.Lock()
			count++
			if err != nil {
				log.Errorf("[%d] couldn't delete droplet %s (%d):\n\t%s", count, droplet.Name, droplet.ID, err)
			} else {
				log.Infof("[%d] deleted %s (%d)", count, droplet.Name, droplet.ID)
			}
			// done with the mutex
			mut.Unlock()

			wg.Done()
		}(v)
	}

	// wait for all the goroutines to finish
	wg.Wait()
}
