package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func createTempFileByKeyAndEmail(key, email string) (string, error) {
	credJsonBody, err := credsFile{
		PrivateKey:  key,
		ClientEmail: email,
	}.serviceAccountJson()
	if err != nil {
		return "", err
	}
	return createTempFile(credJsonBody)
}

func createTempFile(jsonBody string) (string, error) {
	tempDir := os.TempDir()
	f, err := ioutil.TempFile(tempDir, "google-api-go-wrapper.json.")
	if err != nil {
		return "", err
	}

	_, err = f.Write([]byte(jsonBody))
	return f.Name(), err
}

type credsFile struct {
	Type        string `json:"type"`
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
}

func (c credsFile) serviceAccountJson() (string, error) {
	c.Type = "service_account"
	byt, err := json.Marshal(c)
	return string(byt), err
}
