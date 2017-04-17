package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Project struct {
	Id          string `json:"_id,omitempty"`
	Name        string `json:"label,omitempty"`
	GroupId     string `json:"group,omitempty"`
	Description string `json:"description,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
	Files    []*File    `json:"files,omitempty"`

	Notes []*Note                `json:"notes,omitempty"`
	Tags  []string               `json:"tags,omitempty"`
	Info  map[string]interface{} `json:"info,omitempty"`

	Public      *bool         `json:"public,omitempty"`
	Archived    *bool         `json:"archived,omitempty"`
	Permissions []*Permission `json:"permissions,omitempty"`
}

func (c *Client) GetAllProjects() ([]*Project, *http.Response, error) {
	var aerr *Error
	var projects []*Project
	resp, err := c.New().Get("projects").Receive(&projects, &aerr)
	return projects, resp, Coalesce(err, aerr)
}

func (c *Client) GetProject(id string) (*Project, *http.Response, error) {
	var aerr *Error
	var project *Project
	resp, err := c.New().Get("projects/"+id).Receive(&project, &aerr)
	return project, resp, Coalesce(err, aerr)
}

func (c *Client) AddProject(project *Project) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	resp, err := c.New().Post("projects").BodyJSON(project).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

func (c *Client) AddProjectNote(id, text string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	note := &Note{
		Text: text,
	}

	resp, err := c.New().Post("projects/"+id+"/notes").BodyJSON(note).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying project " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) AddProjectTag(id, tag string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	var tagDoc interface{}
	tagDoc = map[string]interface{}{
		"value": tag,
	}

	resp, err := c.New().Post("projects/"+id+"/tags").BodyJSON(tagDoc).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying project " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) ModifyProject(id string, project *Project) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	resp, err := c.New().Put("projects/"+id).BodyJSON(project).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying project " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) DeleteProject(id string) (*http.Response, error) {
	var aerr *Error
	var response *DeletedResponse

	resp, err := c.New().Delete("projects/"+id).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.DeletedCount != 1 {
		return resp, errors.New("Deleting project " + id + " returned " + strconv.Itoa(response.DeletedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) UploadToProject(id string, files ...*UploadSource) (chan int64, chan error) {
	url := "projects/" + id + "/files"
	return c.UploadSimple(url, nil, files...)
}
