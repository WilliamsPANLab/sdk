package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Acquisition struct {
	Id        string `json:"_id,omitempty"`
	Name      string `json:"label,omitempty"`
	SessionId string `json:"session,omitempty"`

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
	Permissions []*Permission `json:"permissions,omitempty"`
}

func (c *Client) GetAllAcquisitions() ([]*Acquisition, *http.Response, error) {
	var aerr *Error
	var acquisitions []*Acquisition
	resp, err := c.New().Get("acquisitions").Receive(&acquisitions, &aerr)
	return acquisitions, resp, Coalesce(err, aerr)
}

func (c *Client) GetAcquisition(id string) (*Acquisition, *http.Response, error) {
	var aerr *Error
	var acquisition *Acquisition
	resp, err := c.New().Get("acquisitions/"+id).Receive(&acquisition, &aerr)
	return acquisition, resp, Coalesce(err, aerr)
}

func (c *Client) AddAcquisition(acquisition *Acquisition) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	resp, err := c.New().Post("acquisitions").BodyJSON(acquisition).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

func (c *Client) AddAcquisitionNote(id, text string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	note := &Note{
		Text: text,
	}

	resp, err := c.New().Post("acquisitions/"+id+"/notes").BodyJSON(note).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying acquisition " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) AddAcquisitionTag(id, tag string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	var tagDoc interface{}
	tagDoc = map[string]interface{}{
		"value": tag,
	}

	resp, err := c.New().Post("acquisitions/"+id+"/tags").BodyJSON(tagDoc).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying acquisition " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) ModifyAcquisition(id string, acquisition *Acquisition) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	resp, err := c.New().Put("acquisitions/"+id).BodyJSON(acquisition).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying acquisition " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) SetAcquisitionInfo(id string, set map[string]interface{}) (*http.Response, error) {
	url := "acquisitions/" + id + "/info"
	return c.setInfo(url, set, false)
}

func (c *Client) ReplaceAcquisitionInfo(id string, replace map[string]interface{}) (*http.Response, error) {
	url := "acquisitions/" + id + "/info"
	return c.replaceInfo(url, replace, false)
}

func (c *Client) DeleteAcquisitionInfoFields(id string, keys []string) (*http.Response, error) {
	url := "acquisitions/" + id + "/info"
	return c.deleteInfoFields(url, keys, false)
}

func (c *Client) DeleteAcquisition(id string) (*http.Response, error) {
	var aerr *Error
	var response *DeletedResponse

	resp, err := c.New().Delete("acquisitions/"+id).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.DeletedCount != 1 {
		return resp, errors.New("Deleting acquisition " + id + " returned " + strconv.Itoa(response.DeletedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) UploadToAcquisition(id string, files ...*UploadSource) (chan int64, chan error) {
	url := "acquisitions/" + id + "/files"
	return c.UploadSimple(url, nil, files...)
}

func (c *Client) ModifyAcquisitionFile(id string, filename string, attributes *FileFields) (*http.Response, *ModifiedAndJobsResponse, error) {
	url := "acquisitions/" + id + "/files/" + filename
	return c.modifyFileAttrs(url, attributes)
}

func (c *Client) SetAcquisitionFileInfo(id string, filename string, set map[string]interface{}) (*http.Response, error) {
	url := "acquisitions/" + id + "/files/" + filename + "/info"
	return c.setInfo(url, set, true)
}

func (c *Client) ReplaceAcquisitionFileInfo(id string, filename string, replace map[string]interface{}) (*http.Response, error) {
	url := "acquisitions/" + id + "/files/" + filename + "/info"
	return c.replaceInfo(url, replace, true)
}

func (c *Client) DeleteAcquisitionFileInfoFields(id string, filename string, keys []string) (*http.Response, error) {
	url := "acquisitions/" + id + "/files/" + filename + "/info"
	return c.deleteInfoFields(url, keys, true)
}

func (c *Client) DownloadFromAcquisition(id string, filename string, destination *DownloadSource) (chan int64, chan error) {
	url := "acquisitions/" + id + "/files/" + filename
	return c.DownloadSimple(url, destination)
}

func (c *Client) GetAcquisitionDownloadUrl(id string, filename string) (string, *http.Response, error) {
	return c.GetTicketDownloadUrl("acquisitions", id, filename)
}

// No progress reporting
func (c *Client) UploadFileToAcquisition(id string, path string) error {
	src := CreateUploadSourceFromFilenames(path)
	progress, result := c.UploadToAcquisition(id, src...)

	// drain and report
	for range progress {
	}
	return <-result
}

// No progress reporting
func (c *Client) DownloadFileFromAcquisition(id, name string, path string) error {
	src := CreateDownloadSourceFromFilename(path)
	progress, result := c.DownloadFromAcquisition(id, name, src)

	// drain and report
	for range progress {
	}
	return <-result
}
