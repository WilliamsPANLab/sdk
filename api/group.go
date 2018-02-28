package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Group struct {
	Id   string `json:"_id,omitempty" bson:"_id"`
	Name string `json:"label,omitempty" bson:"label"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`

	Tags []string `json:"tags,omitempty"`

	// Permissions array is called roles on groups, and groups only
	// https://github.com/scitran/core/issues/662
	Permissions []*Permission `json:"permissions,omitempty"`
}

func (c *Client) GetAllGroups() ([]*Group, *http.Response, error) {
	var aerr *Error
	var groups []*Group
	resp, err := c.New().Get("groups").Receive(&groups, &aerr)
	return groups, resp, Coalesce(err, aerr)
}

func (c *Client) GetGroup(id string) (*Group, *http.Response, error) {
	var aerr *Error
	var group *Group
	resp, err := c.New().Get("groups/"+id).Receive(&group, &aerr)
	return group, resp, Coalesce(err, aerr)
}

func (c *Client) AddGroup(group *Group) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	resp, err := c.New().Post("groups").BodyJSON(group).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

func (c *Client) AddGroupTag(id, tag string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	var tagDoc interface{}
	tagDoc = map[string]interface{}{
		"value": tag,
	}

	resp, err := c.New().Post("groups/"+id+"/tags").BodyJSON(tagDoc).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying group " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) ModifyGroup(id string, group *Group) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	resp, err := c.New().Put("groups/"+id).BodyJSON(group).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying group " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) DeleteGroup(id string) (*http.Response, error) {
	var aerr *Error
	var response *DeletedResponse

	resp, err := c.New().Delete("groups/"+id).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.DeletedCount != 1 {
		return resp, errors.New("Deleting group " + id + " returned " + strconv.Itoa(response.DeletedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}
