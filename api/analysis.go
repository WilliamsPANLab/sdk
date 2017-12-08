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

	Input  bool `json:"input,omitempty"`
	Output bool `json:"output,omitempty"`

	Info map[string]interface{} `json:"info,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
}

type Analysis struct {
	Id     string              `json:"_id,omitempty"`
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

	Files []*AnalysisFile `json:"files,omitempty"`

	Public      bool          `json:"public,omitempty"`
	Permissions []*Permission `json:"permissions,omitempty"`
}

type AnalysisListItem struct {
	Id     string              `json:"_id,omitempty"`
	Name   string              `json:"label,omitempty"`
	Parent *ContainerReference `json:"parent,omitempty"`

	Description string `json:"description,omitempty"`

	// Treat this as a origin of { 'type': 'user', 'id': 'this-field' }
	User string `json:"user,omitempty"`

	Notes []*Note `json:"notes,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`

	Files []*AnalysisFile `json:"files,omitempty"`

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

func (c *Client) AddSessionAnalysis(sessionId string, analysis *Analysis, job *Job) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	body := map[string]interface{}{
		"analysis": analysis,
		"job":      job,
	}

	// 'job=true' flag indicates new behavior
	// https://github.com/scitran/core/issues/908#issuecomment-324766048 item 3
	resp, err := c.New().Post("sessions/"+sessionId+"/analyses?job=true").BodyJSON(body).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

func (c *Client) AddSessionAnalysisNote(sessionId string, analysisId string, text string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	body := &Note{
		Text: text,
	}

	resp, err := c.New().Post("sessions/"+sessionId+"/analyses/"+analysisId+"/notes").BodyJSON(body).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying session analysis on " + sessionId + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) DownloadFromAnalysis(sessionId, analysisId, filename string, destination *DownloadSource) (chan int64, chan error) {
	url := "sessions/" + sessionId + "/analyses/" + analysisId + "/files/" + filename
	return c.DownloadSimple(url, destination)
}

// No progress reporting
func (c *Client) DownloadFileFromAnalysis(sessionId, analysisId, filename, path string) error {
	src := CreateDownloadSourceFromFilename(path)
	progress, result := c.DownloadFromAnalysis(sessionId, analysisId, filename, src)

	// drain and report
	for range progress {
	}
	return <-result
}
