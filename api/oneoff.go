package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Helper func
func (c *Client) modifyFileAttrs(url string, attributes *FileFields) (*http.Response, *ModifiedAndJobsResponse, error) {
	var aerr *Error
	var response *ModifiedAndJobsResponse

	resp, err := c.New().Put(url).BodyJSON(attributes).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, nil, errors.New("Modifying file returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, response, Coalesce(err, aerr)
}

// Helper func
func (c *Client) setInfo(url string, set map[string]interface{}, expectResponse bool) (*http.Response, error) {
	body := map[string]interface{}{
		"set": set,
	}
	return c.postWithOptionalModifiedResponse(url, body, expectResponse)
}

// Helper func
func (c *Client) replaceInfo(url string, replace map[string]interface{}, expectResponse bool) (*http.Response, error) {
	body := map[string]interface{}{
		"replace": replace,
	}
	return c.postWithOptionalModifiedResponse(url, body, expectResponse)
}

// Helper func
func (c *Client) deleteInfoFields(url string, keys []string, expectResponse bool) (*http.Response, error) {
	body := map[string]interface{}{
		"delete": keys,
	}
	return c.postWithOptionalModifiedResponse(url, body, expectResponse)
}

// Helper func
func (c *Client) postWithOptionalModifiedResponse(url string, body interface{}, expectResponse bool) (*http.Response, error) {
	var aerr *Error
	var resp *http.Response
	var err error

	if expectResponse {
		var response *ModifiedResponse
		resp, err = c.New().Post(url).BodyJSON(body).Receive(&response, &aerr)

		if err == nil && aerr == nil && response.ModifiedCount != 1 {
			// Should not have to check this count
			// https://github.com/scitran/core/issues/680
			return resp, errors.New("Modifying file returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
		}
	} else {
		resp, err = c.New().Post(url).BodyJSON(body).Receive(nil, &aerr)
	}

	return resp, Coalesce(err, aerr)
}

// ParseApiKey accepts an API key and returns the hostname, port, key, and any parsing error.
func ParseApiKey(apiKey string) (string, int, string, error) {
	var err error
	host := ""
	port := 443
	key := ""

	splits := strings.Split(apiKey, ":")

	if len(splits) < 2 {
		return host, port, key, errors.New("Invalid API key")
	}

	if len(splits) == 2 {
		host = splits[0]
		key = splits[1]
	} else {
		host = splits[0]
		port, err = strconv.Atoi(splits[1])
		key = splits[len(splits)-1]
	}

	return host, port, key, err
}

// Config represents some of the server's configuration.
type Config struct {
	// Auth holds information about how users can authenticate.
	//
	// NOTE: this can go one layer deeper after multi-auth.
	// https://github.com/scitran/core/pull/652
	Auth map[string]interface{} `json:"auth"`

	// Site holds multi-site registration information.
	// This feature is depreciated.
	Site map[string]interface{} `json:"site"`

	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

func (c *Client) GetConfig() (*Config, *http.Response, error) {
	var aerr *Error
	var config *Config
	resp, err := c.New().Get("config").Receive(&config, &aerr)
	return config, resp, Coalesce(err, aerr)
}

// Version identifies the upgrade level of system components.
type Version struct {
	// Database represents the database schema level.
	Database int `json:"database"`
}

func (c *Client) GetVersion() (*Version, *http.Response, error) {
	var aerr *Error
	var version *Version
	resp, err := c.New().Get("version").Receive(&version, &aerr)
	return version, resp, Coalesce(err, aerr)
}
