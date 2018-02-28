package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

// Key is an API key, saved on a User to access the API.
type Key struct {
	Key string `json:"key"`

	Created  *time.Time `json:"created"`
	LastUsed *time.Time `json:"last_used"`
}

// User represents a single user.
type User struct {
	Id        string `json:"_id,omitempty" bson:"_id"`
	Email     string `json:"email,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	ApiKey    *Key   `json:"api_key,omitempty"`

	Avatar  string            `json:"avatar,omitempty"`
	Avatars map[string]string `json:"avatars,omitempty"`

	Created    *time.Time `json:"created,omitempty"`
	Modified   *time.Time `json:"modified,omitempty"`
	RootAccess *bool      `json:"root,omitempty"`
}

func (c *Client) GetCurrentUser() (*User, *http.Response, error) {
	var aerr *Error
	var user *User
	resp, err := c.New().Get("users/self").Receive(&user, &aerr)
	return user, resp, Coalesce(err, aerr)
}

func (c *Client) GetAllUsers() ([]*User, *http.Response, error) {
	var aerr *Error
	var users []*User
	resp, err := c.New().Get("users").Receive(&users, &aerr)
	return users, resp, Coalesce(err, aerr)
}

func (c *Client) GetUser(id string) (*User, *http.Response, error) {
	var aerr *Error
	var user *User

	resp, err := c.New().Get("users/"+id).Receive(&user, &aerr)
	return user, resp, Coalesce(err, aerr)
}

// AddUser creates a user, and returns the created Id.
func (c *Client) AddUser(user *User) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	resp, err := c.New().Post("users").BodyJSON(user).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

// ModifyUser will update an existing user.
// Only set the fields of user that you want modified.
func (c *Client) ModifyUser(id string, user *User) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	resp, err := c.New().Put("users/"+id).BodyJSON(user).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying user " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

// DeleteUser will delete a user. Returns an error if user was not found or the delete did not succeed.
func (c *Client) DeleteUser(id string) (*http.Response, error) {
	var aerr *Error
	var response *DeletedResponse

	resp, err := c.New().Delete("users/"+id).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.DeletedCount != 1 {
		return resp, errors.New("Deleting user " + id + " returned " + strconv.Itoa(response.DeletedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}
