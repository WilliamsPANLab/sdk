package api

import (
	"net/http"
	"time"
)

// Formula describes a single unit of work.
type Formula struct {
	Inputs  []*Input  `json:"inputs"`
	Target  Target    `json:"target"`
	Outputs []*Output `json:"outputs"`
}

// Input describes an asset that must be present for the formula to execute.
type Input struct {
	Type     string `json:"type,omitempty"`
	URI      string `json:"uri,omitempty"`
	Location string `json:"location,omitempty"`
	VuID     string `json:"vu,omitempty"`
}

// Target describes what the formula will execute.
type Target struct {
	Command []string          `json:"command,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	Dir     string            `json:"dir,omitempty"`
}

// Output describes the creation of an asset after a formula is complete.
type Output struct {
	Type     string `json:"type,omitempty"`
	URI      string `json:"uri,omitempty"`
	Location string `json:"location,omitempty"`
	VuID     string `json:"vu,omitempty"`
}

// Result describes the result of a formula.
type Result struct {
	ExitCode int `json:"exitcode"`
}

// FormulaResult combines a (possibly-modified) Formula with any results.
type FormulaResult struct {
	Formula

	Result Result `json:"result"`
}

// Enum for job states.
type JobState string

const (
	Pending   JobState = "pending"
	Running   JobState = "running"
	Failed    JobState = "failed"
	Complete  JobState = "complete"
	Cancelled JobState = "cancelled"
)

// Enum for job retrieval attempts.
type JobRetrieval int

const (
	JobAquired JobRetrieval = iota
	NoPendingJobs
	JobFailure
)

type JobProfile struct {
	Elapsed int64 `json:"elapsed, omitempty"`
}

type Job struct {
	Id     string `json:"id,omitempty" bson:"_id"`
	GearId string `json:"gear_id,omitempty" bson:"gear_id"`

	State   JobState `json:"state,omitempty"`
	Attempt int      `json:"attempt,omitempty"`
	Origin  *Origin  `json:"origin,omitempty"`

	Config      map[string]interface{} `json:"config,omitempty"`
	Inputs      map[string]interface{} `json:"inputs,omitempty"`
	Destination *ContainerReference    `json:"destination,omitempty"`
	Tags        []string               `json:"tags,omitempty"`

	Request *Formula `json:"request,omitempty"`

	ResultMetadata map[string]interface{} `json:"produced_metadata,omitempty"`
	ResultFiles    []string               `json:"saved_files,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`

	Profile JobProfile `json:"profile,omitempty"`
}

type JobLogStatement struct {
	// The file descriptor the log line came from.
	// Negative values are external logs from Flywheel components;
	// One is standard out;
	// Two is standard err;
	// Other values are invalid.
	FileDescriptor int8 `json:"fd"`

	// The message for this statement. Typically newline-delimited.
	Message string `json:"msg"`
}

type JobLog struct {
	Id   string             `json:"_id,omitempty"`
	Logs []*JobLogStatement `json:"logs,omitempty"`
}

// Get all jobs endpoint is not implemented as it returns a different format
// https://github.com/scitran/core/issues/704

func (c *Client) GetJob(id string) (*Job, *http.Response, error) {
	var aerr *Error
	var job *Job

	resp, err := c.New().Get("jobs/"+id).Receive(&job, &aerr)
	return job, resp, Coalesce(err, aerr)
}

func (c *Client) GetJobLogs(id string) (*JobLog, *http.Response, error) {
	var aerr *Error
	var logs *JobLog

	resp, err := c.New().Get("jobs/"+id+"/logs").Receive(&logs, &aerr)
	return logs, resp, Coalesce(err, aerr)
}

func (c *Client) AddJob(job *Job) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	resp, err := c.New().Post("jobs/add").BodyJSON(job).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

func (c *Client) AddJobLogs(id string, statements []*JobLogStatement) (*http.Response, error) {
	var aerr *Error

	resp, err := c.New().Post("jobs/"+id+"/logs").BodyJSON(statements).Receive(nil, &aerr)
	return resp, Coalesce(err, aerr)
}

func (c *Client) ModifyJob(id string, job *Job) (*http.Response, error) {
	var aerr *Error

	resp, err := c.New().Put("jobs/"+id).BodyJSON(job).Receive(nil, &aerr)
	return resp, Coalesce(err, aerr)
}

func (c *Client) StartNextPendingJob(tags ...string) (JobRetrieval, *Job, *http.Response, error) {
	var aerr *Error
	var job *Job

	params := &struct {
		Tags []string `url:"tags,omitempty"`
	}{
		Tags: tags,
	}

	resp, err := c.New().Get("jobs/next").QueryStruct(params).Receive(&job, &aerr)
	rerr := Coalesce(err, aerr)

	if rerr == nil && job != nil {
		return JobAquired, job, resp, nil
	} else if rerr != nil && resp != nil && resp.StatusCode == 400 {
		return NoPendingJobs, nil, resp, nil
	} else {
		return JobFailure, nil, resp, rerr
	}
}

func (c *Client) HeartbeatJob(id string) (*http.Response, error) {
	var aerr *Error

	// Send empty modification
	empty := map[string]string{}

	resp, err := c.New().Put("jobs/"+id).BodyJSON(empty).Receive(nil, &aerr)
	return resp, Coalesce(err, aerr)
}

func (c *Client) ChangeJobState(id string, state JobState) (*http.Response, error) {
	var aerr *Error
	jobMod := &Job{
		State: state,
	}

	resp, err := c.New().Put("jobs/"+id).BodyJSON(jobMod).Receive(nil, &aerr)
	return resp, Coalesce(err, aerr)
}
