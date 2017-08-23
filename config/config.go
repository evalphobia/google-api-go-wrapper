package config

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

const (
	defaultEnvPrivateKey = "GOOGLE_API_GO_PRIVATEKEY"
	defaultEnvEmail      = "GOOGLE_API_GO_EMAIL"
)

var (
	envEmail      string
	envPrivateKey string
)

func init() {
	envPrivateKey = os.Getenv(defaultEnvPrivateKey)
	envEmail = os.Getenv(defaultEnvEmail)
}

type Config struct {
	// by parameter
	Email      string
	PrivateKey string

	// by file
	Filename string

	Scopes   []string
	TokenURL string
	Timeout  time.Duration
}

func (c Config) Client() (*http.Client, error) {
	conf, err := c.JWTConfig()
	if err != nil {
		return nil, err
	}

	cli := conf.Client(c.NewContext())
	return cli, nil
}

func (c Config) NewContext() context.Context {
	return context.WithValue(oauth2.NoContext, oauth2.HTTPClient, &http.Client{
		Timeout: c.Timeout,
	})
}

func (c Config) JWTConfig() (conf *jwt.Config, err error) {
	switch {
	case c.PrivateKey != "" && c.Email != "":
		conf = newJWTConfigFromParams(c.PrivateKey, c.Email)
	case c.Filename != "":
		conf, err = newJWTConfigFromFilepath(c.Filename)
	case envEmail != "" && envPrivateKey != "":
		conf = newJWTConfigFromParams(envPrivateKey, envEmail)
	default:
		var cred *google.DefaultCredentials
		cred, err = google.FindDefaultCredentials(context.Background(), c.Scopes...)
		if err != nil {
			return nil, err
		}
		if cred.JSON == nil {
			return nil, errors.New("cannot find any environment parameter or required field for google api")
		}
		conf, err = newJWTConfig(cred.JSON)
	}

	if err != nil {
		return nil, err
	}

	conf.Scopes = c.Scopes
	return conf, nil
}

func newJWTConfigFromParams(key, email string) *jwt.Config {
	return &jwt.Config{
		Email:      email,
		PrivateKey: []byte(key),
		TokenURL:   google.JWTTokenURL,
	}
}

func newJWTConfig(jsonKeyData []byte) (*jwt.Config, error) {
	return google.JWTConfigFromJSON(jsonKeyData)
}

func newJWTConfigFromFilepath(path string) (*jwt.Config, error) {
	byt, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return newJWTConfig(byt)
}
