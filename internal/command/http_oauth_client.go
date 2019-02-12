package command

import (
	"crypto/tls"
	"net/http"
)

type AuthClient struct {
	token  string
	Client *http.Client
}

func (d *AuthClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", d.token)
	return d.Client.Do(req)
}

func NewAuthClient(accessToken string, insecureSkipVerify bool) *AuthClient {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecureSkipVerify,
			},
		},
	}

	return &AuthClient{
		token:  accessToken,
		Client: client,
	}
}