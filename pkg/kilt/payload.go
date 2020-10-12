package kilt

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
)



func retrievePayloadViaURL(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	payload, err := ioutil.ReadAll(resp.Body)
	return payload, err
}

func retrievePayloadLocal(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func retrievePayloadBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}
