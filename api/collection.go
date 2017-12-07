package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Collection struct {
	Id          string `json:"_id,omitempty"`
	Name        string `json:"label,omitempty"`
	Curator     string `json:"curator,omitempty"`
	Description string `json:"description,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
	Files    []*File    `json:"files,omitempty"`

	Notes []*Note                `json:"notes,omitempty"`
	Info  map[string]interface{} `json:"info,omitempty"`

	Public      *bool         `json:"public,omitempty"`
	Archived    *bool         `json:"archived,omitempty"`
	Permissions []*Permission `json:"permissions,omitempty"`
}

type collectionNode struct {
	Id    string `json:"_id,omitempty"`
	Level string `json:"level,omitempty"`
}

type collectionOperation struct {
	Operation string            `json:"operation,omitempty"`
	Nodes     []*collectionNode `json:"nodes,omitempty"`
}

type collectionSubmission struct {
	Contents *collectionOperation `json:"contents,omitempty"`
}

func (c *Client) GetAllCollections() ([]*Collection, *http.Response, error) {
	var aerr *Error
	var collections []*Collection
	resp, err := c.New().Get("collections").Receive(&collections, &aerr)
	return collections, resp, Coalesce(err, aerr)
}

func (c *Client) GetCollection(id string) (*Collection, *http.Response, error) {
	var aerr *Error
	var collection *Collection
	resp, err := c.New().Get("collections/"+id).Receive(&collection, &aerr)
	return collection, resp, Coalesce(err, aerr)
}

func (c *Client) GetCollectionSessions(id string) ([]*Session, *http.Response, error) {
	var aerr *Error
	var sessions []*Session
	resp, err := c.New().Get("collections/"+id+"/sessions").Receive(&sessions, &aerr)
	return sessions, resp, Coalesce(err, aerr)
}

func (c *Client) GetCollectionAcquisitions(id string) ([]*Session, *http.Response, error) {
	var aerr *Error
	var sessions []*Session
	resp, err := c.New().Get("collections/"+id+"/acquisitions").Receive(&sessions, &aerr)
	return sessions, resp, Coalesce(err, aerr)
}

func (c *Client) GetCollectionSessionAcquisitions(id string, sid string) ([]*Session, *http.Response, error) {
	var aerr *Error
	var sessions []*Session
	resp, err := c.New().Get("collections/"+id+"/acquisitions?session="+sid).Receive(&sessions, &aerr)
	return sessions, resp, Coalesce(err, aerr)
}

func (c *Client) AddCollection(collection *Collection) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	resp, err := c.New().Post("collections").BodyJSON(collection).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

func (c *Client) addNodesToCollection(id, nodeType string, nodeIds []string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	var nodes []*collectionNode

	for _, acqId := range nodeIds {
		nodes = append(nodes, &collectionNode{
			Id:    acqId,
			Level: nodeType,
		})
	}

	submit := &collectionSubmission{
		Contents: &collectionOperation{
			Operation: "add",
			Nodes:     nodes,
		},
	}

	resp, err := c.New().Put("collections/"+id).BodyJSON(submit).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying collection " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) AddAcquisitionsToCollection(id string, aqids []string) (*http.Response, error) {
	return c.addNodesToCollection(id, "acquisition", aqids)
}

func (c *Client) AddSessionsToCollection(id string, sessionids []string) (*http.Response, error) {
	return c.addNodesToCollection(id, "session", sessionids)
}

func (c *Client) AddCollectionNote(id, text string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	note := &Note{
		Text: text,
	}

	resp, err := c.New().Post("collections/"+id+"/notes").BodyJSON(note).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying collection " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) ModifyCollection(id string, collection *Collection) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	resp, err := c.New().Put("collections/"+id).BodyJSON(collection).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying collection " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) DeleteCollection(id string) (*http.Response, error) {
	var aerr *Error

	// Unlike other delete endpoints, this doesn't return anything. Which is good, but inconsistent.
	// https://github.com/scitran/core/issues/680
	resp, err := c.New().Delete("collections/"+id).Receive(nil, &aerr)

	return resp, Coalesce(err, aerr)
}

func (c *Client) UploadToCollection(id string, files ...*UploadSource) (chan int64, chan error) {
	url := "collections/" + id + "/files"
	return c.UploadSimple(url, nil, files...)
}

func (c *Client) ModifyCollectionFile(id string, filename string, attributes *FileFields) (*http.Response, *ModifiedAndJobsResponse, error) {
	url := "collections/" + id + "/files/" + filename
	return c.modifyFileAttrs(url, attributes)
}

func (c *Client) SetCollectionFileInfo(id string, filename string, set map[string]interface{}) (*http.Response, error) {
	url := "collections/" + id + "/files/" + filename + "/info"
	return c.setInfo(url, set)
}

func (c *Client) ReplaceCollectionFileInfo(id string, filename string, replace map[string]interface{}) (*http.Response, error) {
	url := "collections/" + id + "/files/" + filename + "/info"
	return c.replaceInfo(url, replace)
}

func (c *Client) DeleteCollectionFileInfoFields(id string, filename string, keys []string) (*http.Response, error) {
	url := "collections/" + id + "/files/" + filename + "/info"
	return c.deleteInfoFields(url, keys)
}

func (c *Client) DownloadFromCollection(id string, filename string, destination *DownloadSource) (chan int64, chan error) {
	url := "collections/" + id + "/files/" + filename
	return c.DownloadSimple(url, destination)
}

func (c *Client) GetCollectionDownloadUrl(id string, filename string) (string, *http.Response, error) {
	return c.GetTicketDownloadUrl("collections", id, filename)
}

// No progress reporting
func (c *Client) UploadFileToCollection(id string, path string) error {
	src := CreateUploadSourceFromFilenames(path)
	progress, result := c.UploadToCollection(id, src...)

	// drain and report
	for range progress {
	}
	return <-result
}

// No progress reporting
func (c *Client) DownloadFileFromCollection(id, name string, path string) error {
	src := CreateDownloadSourceFromFilename(path)
	progress, result := c.DownloadFromCollection(id, name, src)

	// drain and report
	for range progress {
	}
	return <-result
}
