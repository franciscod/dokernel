package dokernel

import (
	"bufio"
	"log"
	"os"

	"github.com/franciscod/godo"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func readTokenFromFile() string {
	path, err := homedir.Expand("~/.digitaloceantoken")
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		log.Print("There was an error reading your DigitalOcean token. Make sure you went to https://cloud.digitalocean.com/settings/applications, got a token with Write scope and pasted it into the first line of ~/.digitaloceantoken")
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		return scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ""
}

func clientFromToken(pat string) *godo.Client {
	tokenSrc := &tokenSource{
		AccessToken: pat,
	}

	return godo.NewClient(oauth2.NewClient(oauth2.NoContext, tokenSrc))
}
