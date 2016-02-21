package library

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
	defaultEnvPem        = "GOOGLE_API_GO_PEM"
	defaultEnvEmail      = "GOOGLE_API_GO_EMAIL"
	defaultEnvPrivateKey = "GOOGLE_API_GO_PRIVATEKEY"
)

var (
	envPem        string
	envEmail      string
	envPrivateKey string
)

func init() {
	envPem = os.Getenv(defaultEnvPem)
	envEmail = os.Getenv(defaultEnvEmail)
	envPrivateKey = os.Getenv(defaultEnvPrivateKey)
}

func NewAPIConfig(jsonKey []byte, scope string) (*jwt.Config, error) {
	conf, err := google.JWTConfigFromJSON(jsonKey, scope)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func NewAPIConfigWithPem(pemPath, scope string) (*jwt.Config, error) {
	byt, err := ioutil.ReadFile(pemPath)
	if err != nil {
		return nil, err
	}
	return NewAPIConfig(byt, scope)
}

func NewAPIConfigWithParams(key, email, scope string) *jwt.Config {
	return &jwt.Config{
		Email:      email,
		PrivateKey: []byte(key),
		Scopes:     []string{scope},
		TokenURL:   google.JWTTokenURL,
	}
}

type Config struct {
	Email      string
	PrivateKey []byte
	Scopes     []string
	TokenURL   string
	Timeout    time.Duration
}

func NewAPIConfigWithConfig(c Config) *jwt.Config {
	return &jwt.Config{
		Email:      c.Email,
		PrivateKey: c.PrivateKey,
		Scopes:     c.Scopes,
		TokenURL:   c.TokenURL,
	}
}

func NewAPIConfigWithEnv(scope string) (*jwt.Config, error) {
	switch {
	case envPem != "":
		return NewAPIConfigWithPem(envPem, scope)
	case envEmail != "" && envPrivateKey != "":
		return NewAPIConfigWithParams(envPrivateKey, envEmail, scope), nil
	}
	return nil, errors.New("cannot find any environment parameter for google api")
}

func NewContext(timeout time.Duration) context.Context {
	return context.WithValue(oauth2.NoContext, oauth2.HTTPClient, &http.Client{
		Timeout: timeout,
	})
}
