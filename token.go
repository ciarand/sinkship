package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"strings"
)

// a function that provides a means to retrieve a token
type TokenGetter func() (string, error)

// goes through the provided TokenGetter chain and stops once one reports
// a non-"" value. If any produce errors it'll wrap those up and return 'em.
func getTokenFromChain(getters ...TokenGetter) (string, error) {
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

// opens the "token" file and reads it into a string
func getTokenFromFile() (string, error) {
	bytes, err := ioutil.ReadFile("token")
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// checks the DO_TOKEN env variable
func getTokenFromEnv() (string, error) {
	return os.Getenv("DO_TOKEN"), nil
}

// checks the "-token" flag on the CLI
func getTokenFromCli() (string, error) {
	var str *string

	flag.StringVar(str, "token", "", "The token to use with the DO API")
	flag.Parse()

	return *str, nil
}
