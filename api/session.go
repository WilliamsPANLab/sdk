package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Subject struct {
	Id   string `json:"_id,omitempty"`
	Code string `json:"code,omitempty"`

	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`

	Sex  string                 `json:"sex,omitempty"`
	Age  int                    `json:"age,omitempty"`
	Info map[string]interface{} `json:"info,omitempty"`
}

type Session struct {
	Id        string `json:"_id,omitempty"`
	Name      string `json:"label,omitempty"`
	GroupId   string `json:"group,omitempty"`
	ProjectId string `json:"project,omitempty"`

	Subject   *Subject   `json:"subject,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
	Timezone  string     `json:"timezone,omitempty"`
	Uid       string     `json:"uid,omitempty"`

	Notes []*Note                `json:"notes,omitempty"`
	Tags  []string               `json:"tags,omitempty"`
	Info  map[string]interface{} `json:"info,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
	Files    []*File    `json:"files,omitempty"`

	Public      bool          `json:"public,omitempty"`
	Archived    *bool         `json:"archived,omitempty"`
	Permissions []*Permission `json:"permissions,omitempty"`

	Analyses []*Analysis `json:"analyses,omitempty"`
}

func (c *Client) GetAllSessions() ([]*Session, *http.Response, error) {
	var aerr *Error
	var sessions []*Session
	resp, err := c.New().Get("sessions").Receive(&sessions, &aerr)
	return sessions, resp, Coalesce(err, aerr)
}

func (c *Client) GetSession(id string) (*Session, *http.Response, error) {
	var aerr *Error
	var session *Session
	resp, err := c.New().Get("sessions/"+id).Receive(&session, &aerr)
	return session, resp, Coalesce(err, aerr)
}
func (c *Client) GetSessionAcquisitions(id string) ([]*Acquisition, *http.Response, error) {
	var aerr *Error
	var acquisitions []*Acquisition
	resp, err := c.New().Get("sessions/"+id+"/acquisitions").Receive(&acquisitions, &aerr)
	return acquisitions, resp, Coalesce(err, aerr)
}

func (c *Client) AddSession(session *Session) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	resp, err := c.New().Post("sessions").BodyJSON(session).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

func (c *Client) AddSessionNote(id, text string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	note := &Note{
		Text: text,
	}

	resp, err := c.New().Post("sessions/"+id+"/notes").BodyJSON(note).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying session " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) AddSessionTag(id, tag string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	var tagDoc interface{}
	tagDoc = map[string]interface{}{
		"value": tag,
	}

	resp, err := c.New().Post("sessions/"+id+"/tags").BodyJSON(tagDoc).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying session " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) ModifySession(id string, session *Session) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	resp, err := c.New().Put("sessions/"+id).BodyJSON(session).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying session " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) SetSessionInfo(id string, set map[string]interface{}) (*http.Response, error) {
	url := "sessions/" + id + "/info"
	return c.setInfo(url, set, false)
}

func (c *Client) ReplaceSessionInfo(id string, replace map[string]interface{}) (*http.Response, error) {
	url := "sessions/" + id + "/info"
	return c.replaceInfo(url, replace, false)
}

func (c *Client) DeleteSessionInfoFields(id string, keys []string) (*http.Response, error) {
	url := "sessions/" + id + "/info"
	return c.deleteInfoFields(url, keys, false)
}

func (c *Client) DeleteSession(id string) (*http.Response, error) {
	var aerr *Error
	var response *DeletedResponse

	resp, err := c.New().Delete("sessions/"+id).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.DeletedCount != 1 {
		return resp, errors.New("Deleting session " + id + " returned " + strconv.Itoa(response.DeletedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) UploadToSession(id string, files ...*UploadSource) (chan int64, chan error) {
	url := "sessions/" + id + "/files"
	return c.UploadSimple(url, nil, files...)
}

func (c *Client) ModifySessionFile(id string, filename string, attributes *FileFields) (*http.Response, *ModifiedAndJobsResponse, error) {
	url := "sessions/" + id + "/files/" + filename
	return c.modifyFileAttrs(url, attributes)
}

func (c *Client) SetSessionFileInfo(id string, filename string, set map[string]interface{}) (*http.Response, error) {
	url := "sessions/" + id + "/files/" + filename + "/info"
	return c.setInfo(url, set, true)
}

func (c *Client) ReplaceSessionFileInfo(id string, filename string, replace map[string]interface{}) (*http.Response, error) {
	url := "sessions/" + id + "/files/" + filename + "/info"
	return c.replaceInfo(url, replace, true)
}

func (c *Client) DeleteSessionFileInfoFields(id string, filename string, keys []string) (*http.Response, error) {
	url := "sessions/" + id + "/files/" + filename + "/info"
	return c.deleteInfoFields(url, keys, true)
}

func (c *Client) DownloadFromSession(id string, filename string, destination *DownloadSource) (chan int64, chan error) {
	url := "sessions/" + id + "/files/" + filename
	return c.DownloadSimple(url, destination)
}

func (c *Client) GetSessionDownloadUrl(id string, filename string) (string, *http.Response, error) {
	return c.GetTicketDownloadUrl("sessions", id, filename)
}

// No progress reporting
func (c *Client) UploadFileToSession(id string, path string) error {
	src := CreateUploadSourceFromFilenames(path)
	progress, result := c.UploadToSession(id, src...)

	// drain and report
	for range progress {
	}
	return <-result
}

// No progress reporting
func (c *Client) DownloadFileFromSession(id, name string, path string) error {
	src := CreateDownloadSourceFromFilename(path)
	progress, result := c.DownloadFromSession(id, name, src)

	// drain and report
	for range progress {
	}
	return <-result
}
