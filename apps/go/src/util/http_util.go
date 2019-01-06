package util

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	logging "github.com/op/go-logging"
)

var cacheDir = ""

func init() {
	tmpdir := os.Getenv("TMP")
	if tmpdir == "" {
		return
	}
	cacheDir = filepath.Join(tmpdir, "Cobalt")
	os.MkdirAll(cacheDir, os.ModePerm)
}

// GetURL is a wrapper around common functionality to get blah
func GetURL(url string, log *logging.Logger) ([]byte, error) {
	// Compute the hash of the url so we stop hitting the internet so much
	var urlHash string
	hasher := sha1.New()
	hasher.Write([]byte(url))
	urlHash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	cacheFile := filepath.Join(cacheDir, urlHash+".cached")

	// Try to read the cached content.  If it's there, great!
	result, err := ioutil.ReadFile(cacheFile)
	if err == nil {
		log.Debugf("returning cached %s from %s", url, urlHash)
		return result, err
	}

	// Not there.  Fetch it.
	log.Debugf("fetching %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)

	// If we got it, write back to the cache!
	if err == nil {
		log.Debugf("writing back to %s", cacheFile)
		ioutil.WriteFile(cacheFile, result, 0644)
	}

	// return the original content and error
	return result, nil
}

// GetJSON is useful for grabbing remote, untyped json
func GetJSON(url string, log *logging.Logger) (result map[string]*json.RawMessage, err error) {
	bytes, err := GetURL(url, log)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, result)
	return
}

/*
func ParseRawMessage(json map[string]*json.RawMessage, path ...string) (*json.RawMessage, error) {
	leaf := json
	count := len(path) - 1
	for i := 0; ; i++ {
		name := path[i]
		next := json[name]
		if next == nil {
			return nil, fmt.Errorf("could not find key \"%s\" while descending path \"%v\" in json", name, path)
		}
		if i == count {
			return next, nil
		}
		err := json.Unmarshal(*next, json)
		if err != nil {
			return nil, fmt.Errorf("error parsing json key \"%s\" while descending path \"%v\": %v", name, path, err)
		}
	}
}
*/
