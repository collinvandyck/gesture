// convenience http client methods
package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	// an http client that is re-used
	httpClient = http.Client{}
)

// checks the header of a particular url to see if it's  equal to a response code
func ResponseHeaderHasCode(url string, code int) (bool, error) {
	resp, err := httpClient.Head(url) // will follow redirects
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return (resp.StatusCode == code), nil
}

func ResponseHeaderContentType(url string) (string, error) {
	resp, err := httpClient.Head(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return resp.Header.Get("Content-Type"), nil
	}
	return "", errors.New("Non-OK status code")
}

func ResolveRedirects(url string) (string, error) {
	resp, err := httpClient.Head(url) // will follow redirects
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // not sure if i have to do this with a head response
	expanded := resp.Request.URL.String()
	if expanded != url {
		return expanded, nil
	}
	return "", nil
}

// GETs a url and returns its body as a []byte
func GetUrl(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Bad response code: %d", resp.StatusCode))
	}
	return body, nil
}
