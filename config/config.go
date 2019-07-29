package config

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

const (
	// see https://godoc.org/golang.org/x/oauth2/google#FindDefaultCredentials
	defaultEnvCredFile = "GOOGLE_APPLICATION_CREDENTIALS"

	defaultEnvPrivateKey = "GOOGLE_API_GO_PRIVATEKEY"
	defaultEnvEmail      = "GOOGLE_API_GO_EMAIL"
	defaultEnvJSON       = "GOOGLE_API_GO_JSON"
	defaultEnvUseIAMRole = "GOOGLE_API_GO_USE_IAMROLE"

	// for oauth token creds
	defaultEnvOAuthClientID     = "GOOGLE_API_OAUTH_CLIENT_ID"
	defaultEnvOAuthClientSecret = "GOOGLE_API_OAUTH_CLIENT_SECRET"
	defaultEnvOAuthRedirectURL  = "GOOGLE_API_OAUTH_REDIRECT_URL"
	defaultEnvOAuthCode         = "GOOGLE_API_OAUTH_CODE"
	defaultEnvOAuthTokenFile    = "GOOGLE_API_OAUTH_TOKEN_FILE"
)

var (
	envCredFile          string
	envEmail             string
	envPrivateKey        string
	envJSON              string
	envUseIAMRole        bool
	envOAuthClientID     string
	envOAuthClientSecret string
	envOAuthRedirectURL  string
	envOAuthCode         string
	envOAuthTokenFile    string
)

func init() {
	envCredFile = os.Getenv(defaultEnvCredFile)
	envPrivateKey = os.Getenv(defaultEnvPrivateKey)
	envEmail = os.Getenv(defaultEnvEmail)
	envJSON = os.Getenv(defaultEnvJSON)
	envUseIAMRole, _ = strconv.ParseBool(os.Getenv(defaultEnvUseIAMRole))
	envOAuthClientID = os.Getenv(defaultEnvOAuthClientID)
	envOAuthClientSecret = os.Getenv(defaultEnvOAuthClientSecret)
	envOAuthRedirectURL = os.Getenv(defaultEnvOAuthRedirectURL)
	envOAuthCode = os.Getenv(defaultEnvOAuthCode)
	envOAuthTokenFile = os.Getenv(defaultEnvOAuthTokenFile)
}

type Config struct {
	// by parameter
	Email      string
	PrivateKey string

	// by file
	Filename string

	// by OAuth setting
	OAuthClientID     string
	OAuthClientSecret string
	OAuthRefreshToken string
	OAuthRedirectURL  string
	OAuthCode         string
	OAuthTokenFile    string

	Scopes   []string
	TokenURL string
	Timeout  time.Duration

	CredsJSONBody    string
	UseTempCredsFile bool
	// tempCredsFilePath is filled by CredsFilePath when UseTempCredsFile is true.
	tempCredsFilePath string

	UseIAMRole bool
}

func (c Config) Client() (*http.Client, error) {
	if c.UseIAMRole || envUseIAMRole {
		return google.DefaultClient(c.NewContext())
	}
	if c.hasOAuthClient() {
		return c.NewOAuthClient()
	}

	conf, err := c.JWTConfig()
	if err != nil {
		return nil, err
	}

	cli := conf.Client(c.NewContext())
	return cli, nil
}

// CredsFilePath returns credential file path.
// if UseTempCredsFile is true, then temporary creds json file will be created.
func (c *Config) CredsFilePath() (string, error) {
	switch {
	case c.Filename != "":
		return c.Filename, nil
	case c.UseTempCredsFile:
		// create temporary creds file
		var err error
		switch {
		case c.PrivateKey != "" && c.Email != "":
			c.tempCredsFilePath, err = createTempFileByKeyAndEmail(c.PrivateKey, c.Email)
			return c.tempCredsFilePath, err
		case c.CredsJSONBody != "":
			c.tempCredsFilePath, err = createTempFile(c.CredsJSONBody)
			return c.tempCredsFilePath, err
		case envPrivateKey != "" && envEmail != "":
			c.tempCredsFilePath, err = createTempFileByKeyAndEmail(envPrivateKey, envEmail)
			return c.tempCredsFilePath, err
		case envJSON != "":
			c.tempCredsFilePath, err = createTempFile(envJSON)
			return c.tempCredsFilePath, err
		}
	case envCredFile != "":
		return envCredFile, nil
	}

	return "", errors.New("Cannot find google creds file path")
}

// DeleteTempCredsFile deletes temporary creds json file.
func (c Config) DeleteTempCredsFile() error {
	if c.tempCredsFilePath == "" {
		return nil
	}
	return os.Remove(c.tempCredsFilePath)
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
	case c.CredsJSONBody != "":
		conf, err = newJWTConfig([]byte(c.CredsJSONBody))
	case c.Filename != "":
		conf, err = newJWTConfigFromFilepath(c.Filename)
	case envEmail != "" && envPrivateKey != "":
		conf = newJWTConfigFromParams(envPrivateKey, envEmail)
	case envJSON != "":
		conf, err = newJWTConfig([]byte(envJSON))
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

func newJWTConfig(jsonKeyData []byte) (*jwt.Config, error) {
	return google.JWTConfigFromJSON(jsonKeyData)
}

func newJWTConfigFromParams(key, email string) *jwt.Config {
	return &jwt.Config{
		Email:      email,
		PrivateKey: []byte(key),
		TokenURL:   google.JWTTokenURL,
	}
}

func newJWTConfigFromFilepath(path string) (*jwt.Config, error) {
	byt, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return newJWTConfig(byt)
}
