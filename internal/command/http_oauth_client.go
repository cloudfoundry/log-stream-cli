package command

import (
	"crypto/tls"
	"net/http"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type AuthDoer struct {
	token  string
	Client Doer
}

func (d *AuthDoer) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", d.token)
	return d.Client.Do(req)
}

func NewDoer(accessToken string, insecureSkipVerify bool) Doer {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecureSkipVerify,
			},
		},
	}

	return &AuthDoer{
		token:  accessToken,
		Client: client,
	}
}