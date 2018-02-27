package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

// Two extra fields beyond File
// https://github.com/scitran/core/issues/908#issuecomment-324766048 item 1
type AnalysisFile struct {
	Name   string  `json:"name,omitempty"`
	Origin *Origin `json:"origin,omitempty"`
	Size   int     `json:"size,omitempty"`

	Modality     string   `json:"modality,omitempty"`
	Mimetype     string   `json:"mimetype,omitempty"`
	Measurements []string `json:"measurements,omitempty"`
	Type         string   `json:"type,omitempty"`
	Tags         []string `json:"tags,omitempty"`

	Info map[string]interface{} `json:"info,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
}

type AdhocAnalysis struct {
	Name        string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`

	// Treat this as a origin of { 'type': 'user', 'id': 'this-field' }
	User string `json:"user,omitempty"`

	Notes []*Note `json:"notes,omitempty"`

	Inputs []*FileReference `json:"inputs,omitempty"`
}

type Analysis struct {
	Id     string              `json:"_id,omitempty" bson:"_id"`
	Name   string              `json:"label,omitempty"`
	Parent *ContainerReference `json:"parent,omitempty"`

	Description string `json:"description,omitempty"`

	// Treat this as a origin of { 'type': 'user', 'id': 'this-field' }
	User string `json:"user,omitempty"`

	Notes []*Note `json:"notes,omitempty"`

	// For now, jobs are always inflated by the endpoints we fetch them through.
	// https://github.com/scitran/core/issues/908#issuecomment-324766048 item 2
	Job *Job `json:"job,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`

	Inputs []*AnalysisFile `json:"inputs,omitempty"`
	Files  []*AnalysisFile `json:"files,omitempty"`

	Public      bool          `json:"public,omitempty"`
	Permissions []*Permission `json:"permissions,omitempty"`

	Tags []string               `json:"tags,omitempty"`
	Info map[string]interface{} `json:"info,omitempty"`
}

type AnalysisListItem struct {
	Id     string              `json:"_id,omitempty" bson:"_id"`
	Name   string              `json:"label,omitempty"`
	Parent *ContainerReference `json:"parent,omitempty"`

	Description string `json:"description,omitempty"`

	// Treat this as a origin of { 'type': 'user', 'id': 'this-field' }
	User string `json:"user,omitempty"`

	Notes []*Note `json:"notes,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`

	Inputs []*AnalysisFile `json:"inputs,omitempty"`
	Files  []*AnalysisFile `json:"files,omitempty"`

	Public      bool          `json:"public,omitempty"`
	Permissions []*Permission `json:"permissions,omitempty"`

	JobId string `json:"job,omitempty"`
}

// This may be cleaner if we had an abstract container struct
func (c *Client) GetAnalyses(cont_name string, cid string, sub_cont string) ([]*AnalysisListItem, *http.Response, error) {
	var aerr *Error
	var analyses []*AnalysisListItem
	var url string

	// Check to see if sub_cont is an empty string
	if sub_cont == "" {
		url = cont_name + "/" + cid + "/analyses"
	} else {
		url = cont_name + "/" + cid + "/" + sub_cont + "/analyses"
	}

	resp, err := c.New().Get(url).Receive(&analyses, &aerr)
	return analyses, resp, Coalesce(err, aerr)
}

func (c *Client) GetAnalysis(id string) (*Analysis, *http.Response, error) {
	var aerr *Error
	var analysis *Analysis

	// inflate_job flag is set to avoid a dynamic type on Analysis.Job
	resp, err := c.New().Get("analyses/"+id+"?inflate_job=true").Receive(&analysis, &aerr)
	return analysis, resp, Coalesce(err, aerr)
}

func (c *Client) AddAnalysisNote(analysisId string, text string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	body := &Note{
		Text: text,
	}

	resp, err := c.New().Post("analyses/"+analysisId+"/notes").BodyJSON(body).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying analysis " + analysisId + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) AddAnalysisTag(id, tag string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	var tagDoc interface{}
	tagDoc = map[string]interface{}{
		"value": tag,
	}

	resp, err := c.New().Post("analyses/"+id+"/tags").BodyJSON(tagDoc).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying analysis" + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) SetAnalysisInfo(id string, set map[string]interface{}) (*http.Response, error) {
	url := "analyses/" + id + "/info"
	return c.setInfo(url, set, false)
}

func (c *Client) ReplaceAnalysisInfo(id string, replace map[string]interface{}) (*http.Response, error) {
	url := "analyses/" + id + "/info"
	return c.replaceInfo(url, replace, false)
}

func (c *Client) DeleteAnalysisInfoFields(id string, keys []string) (*http.Response, error) {
	url := "analyses/" + id + "/info"
	return c.deleteInfoFields(url, keys, false)
}

func (c *Client) UploadToAnalysis(id string, files ...*UploadSource) (chan int64, chan error) {
	url := "analyses/" + id + "/files"
	return c.UploadSimple(url, nil, files...)
}

func (c *Client) DownloadFromAnalysis(analysisId, filename string, destination *DownloadSource) (chan int64, chan error) {
	url := "analyses/" + analysisId + "/files/" + filename
	return c.DownloadSimple(url, destination)
}

func (c *Client) DownloadInputFromAnalysis(analysisId, filename string, destination *DownloadSource) (chan int64, chan error) {
	url := "analyses/" + analysisId + "/inputs/" + filename
	return c.DownloadSimple(url, destination)
}

func (c *Client) GetAnalysisDownloadUrl(id string, filename string) (string, *http.Response, error) {
	return c.GetTicketDownloadUrl("analyses", id, filename)
}

func (c *Client) GetAnalysisInputDownloadUrl(id string, filename string) (string, *http.Response, error) {
	url := "analyses/" + id + "/inputs/" + filename
	return c.GetTicketDownloadUrlFromUrl(url)
}

// No progress reporting
func (c *Client) UploadFileToAnalysis(id string, path string) error {
	return c.UploadFilesToAnalysis(id, []string{path})
}

// No progress reporting
func (c *Client) UploadFilesToAnalysis(id string, paths []string) error {
	src := CreateUploadSourceFromFilenames(paths...)
	progress, result := c.UploadToAnalysis(id, src...)

	// drain and report
	for range progress {
	}
	return <-result
}

// No progress reporting
func (c *Client) DownloadFileFromAnalysis(analysisId, filename, path string) error {
	src := CreateDownloadSourceFromFilename(path)
	progress, result := c.DownloadFromAnalysis(analysisId, filename, src)

	// drain and report
	for range progress {
	}
	return <-result
}

// No progress reporting
func (c *Client) DownloadInputFileFromAnalysis(analysisId, filename, path string) error {
	src := CreateDownloadSourceFromFilename(path)
	progress, result := c.DownloadInputFromAnalysis(analysisId, filename, src)

	// drain and report
	for range progress {
	}
	return <-result
}
