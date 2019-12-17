package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	defaultOAuthTokenFile = "oauth_token.json"
)

// NewOAuthClient creates http.Client from OAuth parameters.
func (c Config) NewOAuthClient() (*http.Client, error) {
	conf, err := c.oauthConfig()
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, &http.Client{
		Timeout: c.Timeout,
	})

	// check existence of oauth token file.
	tokenFile := c.getOAuthTokenFile()
	if _, err := os.Stat(tokenFile); err != nil {
		// or create oauth token file with code.
		code := c.getOAuthCode()
		tok, err := conf.Exchange(ctx, code)
		if err != nil {
			return nil, err
		}

		err = c.saveOAuthTokenFile(tok)
		if err != nil {
			return nil, err
		}
	}

	tok, err := c.loadOAuthTokenFromFile(tokenFile)
	if err != nil {
		return nil, err
	}

	ts := conf.TokenSource(ctx, tok)
	_, err = ts.Token()
	if err != nil {
		return nil, err
	}

	return oauth2.NewClient(ctx, ts), nil
}

// GetOAuthCodeURL returns URL to get oauth code.
func (c Config) GetOAuthCodeURL() string {
	conf, _ := c.oauthConfig()
	return conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (c Config) oauthConfig() (*oauth2.Config, error) {
	if c.getOAuthCredentials() == "" {
		return &oauth2.Config{
			ClientID:     c.getOAuthClientID(),
			ClientSecret: c.getOAuthClientSecret(),
			RedirectURL:  c.getOAuthRedirectURL(),
			Scopes:       c.Scopes,
			Endpoint:     google.Endpoint,
		}, nil
	}

	b, err := ioutil.ReadFile(c.getOAuthCredentials())
	if err != nil {
		return nil, err
	}
	return google.ConfigFromJSON(b)
}

// saveOAuthTokenFile saves oauth token json on local file.
func (c Config) saveOAuthTokenFile(token *oauth2.Token) error {
	jsonText, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.getOAuthTokenFile(), jsonText, 0777)
}

func (c Config) loadOAuthTokenFromFile(path string) (*oauth2.Token, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("Cannot find oauth token file: [%s]", path)
	}

	byt, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	tok := new(oauth2.Token)
	err = json.Unmarshal(byt, tok)
	return tok, err
}

func (c Config) useOAuthClient() bool {
	if c.NoOAuthClient {
		return false
	}
	switch {
	case c.getOAuthCredentials() != "",
		c.getOAuthClientID() != "":
		return true
	}
	return false
}

func (c Config) getOAuthCredentials() string {
	if c.OAuthCredsFile != "" {
		return c.OAuthCredsFile
	}
	return envOAuthCredsFile
}

func (c Config) getOAuthClientID() string {
	if c.OAuthClientID != "" {
		return c.OAuthClientID
	}
	return envOAuthClientID
}

func (c Config) getOAuthClientSecret() string {
	if c.OAuthClientSecret != "" {
		return c.OAuthClientSecret
	}
	return envOAuthClientSecret
}

func (c Config) getOAuthRedirectURL() string {
	if c.OAuthRedirectURL != "" {
		return c.OAuthRedirectURL
	}
	return envOAuthRedirectURL
}

func (c Config) getOAuthCode() string {
	if c.OAuthCode != "" {
		return c.OAuthCode
	}
	return envOAuthCode
}

func (c Config) getOAuthTokenFile() string {
	switch {
	case c.OAuthTokenFile != "":
		return c.OAuthTokenFile
	case envOAuthTokenFile != "":
		return envOAuthTokenFile
	}
	return defaultOAuthTokenFile
}
