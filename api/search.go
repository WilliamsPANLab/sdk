package api

import (
	"net/http"
	"strconv"
)

// Enum for Return Types.
type SearchType string

const (
	FileString        SearchType = "file"
	AcquisitionString SearchType = "acquisition"
	SessionString     SearchType = "session"
	AnalysisString    SearchType = "analysis"
	CollectionString  SearchType = "collection"
)

type SearchQuery struct {
	// ReturnType sets the type of search results.
	ReturnType SearchType `json:"return_type"`

	// SearchString represents the search query.
	SearchString string `json:"search_string,omitempty"`

	// Limit determines the maximum number of search results.
	// Set to a negative value to return all results. Set to zero to use a default limit.
	Limit int `json:"limit,omitempty"`

	// IncludeInaccessible, when set, will include data that the current user does not have access to read.
	IncludeInaccessible bool `json:"all_data,omitempty"`

	// Filters is a set of ElasticSearch filters to use in the search.
	// https://www.elastic.co/guide/en/elasticsearch/reference/current/term-level-queries.html
	Filters []interface{} `json:"filters,omitempty"`
}

type ProjectSearchResponse struct {
	Id   string `json:"_id,omitempty"`
	Name string `json:"label,omitempty"`
}
type GroupSearchResponse struct {
	Id   string `json:"_id,omitempty"`
	Name string `json:"label,omitempty"`
}

// Timestamp fields should be time.Time
// https://github.com/flywheel-io/sdk/issues/52

type SessionSearchResponse struct {
	Id        string `json:"_id,omitempty"`
	Archived  bool   `json:"archived,omitempty"`
	Name      string `json:"label,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Created   string `json:"created,omitempty"`
}
type AcquisitionSearchResponse struct {
	Id        string `json:"_id,omitempty"`
	Archived  bool   `json:"archived,omitempty"`
	Name      string `json:"label,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Created   string `json:"created,omitempty"`
}
type SubjectSearchResponse struct {
	Code string `json:"code,omitempty"`
}
type FileSearchResponse struct {
	Measurements []string `json:"measurements,omitempty"`
	Created      string   `json:"created,omitempty"`
	Type         string   `json:"type,omitempty"`
	Name         string   `json:"name,omitempty"`
	Size         int      `json:"size,omitempty"`
}
type AnalysisSearchResponse struct {
	Id      string `json:"_id,omitempty"`
	Name    string `json:"label,omitempty"`
	User    string `json:"user,omitempty"`
	Created string `json:created,omitempty"`
}
type ParentSearchResponse struct {
	Type string `json:"type,omitempty"`
	Id   string `json:"_id,omitempty"`
}
type CollectionSearchResponse struct {
	Id      string `json:"_id,omitempty"`
	Name    string `json:"label,omitempty"`
	Curator string `json:"curator,omitempty"`
	Created string `json:"created,omitempty"`
}

// SearchResponse for the SearchResponse
type SearchResponse struct {
	Project     *ProjectSearchResponse     `json:"project,omitempty"`
	Group       *GroupSearchResponse       `json:"group,omitempty"`
	Session     *SessionSearchResponse     `json:"session,omitempty"`
	Acquisition *AcquisitionSearchResponse `json:"acquisition,omitempty"`
	Subject     *SubjectSearchResponse     `json:"subject,omitempty"`
	File        *FileSearchResponse        `json:"file,omitempty"`
	Collection  *CollectionSearchResponse  `json:"collection,omitempty"`
	Permissions []*Permission              `json:"permissions,omitempty"`
	Analysis    *AnalysisSearchResponse    `json:"analysis,omitempty"`
	Parent      *ParentSearchResponse      `json:"parent,omitempty"`
}

// Search runs a query, returning up to limit results.
func (c *Client) Search(search_query *SearchQuery) ([]*SearchResponse, *http.Response, error) {
	var aerr *Error
	var response []*SearchResponse

	if search_query.Limit == 0 {
		search_query.Limit = 100
	}

	url := "dataexplorer/search?simple=true&size="

	if search_query.Limit >= 0 {
		url += strconv.Itoa(search_query.Limit)
	} else {
		url += "all"
	}

	resp, err := c.New().Post(url).BodyJSON(search_query).Receive(&response, &aerr)

	return response, resp, Coalesce(err, aerr)
}

// SearchResponse is used for endpoints of data_explorer
type RawSearchResponse struct {
	Id     string          `json:"_id"`
	Source *SearchResponse `json:"_source,omitempty"`
}

// Because the endpoint returns a key results which is a list of responses
type RawSearchResponseList struct {
	Results []*RawSearchResponse `json:"results,omitempty"`
}

// SearchRaw is left in for compatibility reasons. You should probably use Search.
func (c *Client) SearchRaw(search_query *SearchQuery) (*RawSearchResponseList, *http.Response, error) {
	var aerr *Error
	var response *RawSearchResponseList

	resp, err := c.New().Post("dataexplorer/search").BodyJSON(search_query).Receive(&response, &aerr)

	return response, resp, Coalesce(err, aerr)
}
