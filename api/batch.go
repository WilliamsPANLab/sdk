package api

import (
	"net/http"
	"time"
)

type Batch struct {
	Id     string `json:"_id,omitempty" bson:"_id"`
	GearId string `json:"gear_id,omitempty"`

	State  JobState `json:"state,omitempty"`
	Origin *Origin  `json:"origin,omitempty"`

	Config map[string]interface{} `json:"config,omitempty"`
	JobIds []string               `json:"jobs,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
}

type BatchProposal struct {
	Id     string                 `json:"_id,omitempty" bson:"_id"`
	GearId string                 `json:"gear_id,omitempty"`
	Config map[string]interface{} `json:"config,omitempty"`
	State  string                 `json:"state,omitempty"`
	Origin *Origin                `json:"origin,omitempty"`

	Proposal interface{} `json:"proposal,omitempty"`

	Ambiguous          []interface{} `json:"ambiguous,omitempty"`
	MissingPermissions []interface{} `json:"improper_permissions,omitempty"`
	Matched            []interface{} `json:"matched,omitempty"`
	NotMatched         []interface{} `json:"not_matched,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
}

func (c *Client) GetAllBatches() ([]*Batch, *http.Response, error) {
	var aerr *Error
	var batchs []*Batch
	resp, err := c.New().Get("batch").Receive(&batchs, &aerr)
	return batchs, resp, Coalesce(err, aerr)
}

func (c *Client) GetBatch(id string) (*Batch, *http.Response, error) {
	var aerr *Error
	var batch *Batch
	resp, err := c.New().Get("batch/"+id).Receive(&batch, &aerr)
	return batch, resp, Coalesce(err, aerr)
}

func (c *Client) ProposeBatch(gearId string, config map[string]interface{}, tags []string, targets []*ContainerReference) (*BatchProposal, *http.Response, error) {
	var aerr *Error
	var proposal *BatchProposal

	batch := &struct {
		GearId  string                 `json:"gear_id"`
		Config  map[string]interface{} `json:"config"`
		Tags    []string               `json:"tags"`
		Targets []*ContainerReference  `json:"targets"`
	}{
		GearId:  gearId,
		Config:  config,
		Tags:    tags,
		Targets: targets,
	}

	resp, err := c.New().Post("batch").BodyJSON(batch).Receive(&proposal, &aerr)
	return proposal, resp, Coalesce(err, aerr)
}

func (c *Client) StartBatch(id string) ([]*Job, *http.Response, error) {
	var aerr *Error
	var jobs []*Job

	resp, err := c.New().Post("batch/"+id+"/run").Receive(&jobs, &aerr)
	return jobs, resp, Coalesce(err, aerr)
}

func (c *Client) CancelBatch(id string) (int, *http.Response, error) {
	var aerr *Error
	response := &struct {
		Cancelled int `json:"number_cancelled"`
	}{
		Cancelled: 0,
	}

	resp, err := c.New().Post("batch/"+id+"/cancel").Receive(&response, &aerr)
	return response.Cancelled, resp, Coalesce(err, aerr)
}
