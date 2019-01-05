package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// GetURL is a wrapper around common functionality to get blah
func GetURL(url string) (result []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)
	return
}

// GetJSON is useful for grabbing remote, untyped json
func GetJSON(url string) (result map[string]interface{}, err error) {
	bytes, err := GetURL(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, result)
	return
}
