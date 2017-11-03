package api

import (
	"net/http"
	"time"
)

type GearSource struct {
	Commit     string `json:"git-commit,omitempty"`
	RootfsHash string `json:"rootfs-hash,omitempty"`
	RootfsUrl  string `json:"rootfs-url,omitempty"`
}

type Gear struct {
	Name        string `json:"name,omitempty"`
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`

	Version  string `json:"version,omitempty"`
	Flywheel string `json:"flywheel,omitempty"`

	Inputs map[string]map[string]interface{} `json:"inputs"`
	Config map[string]map[string]interface{} `json:"config"`

	Environment map[string]string `json:"environment,omitempty"`
	Command     string            `json:"command,omitempty"`

	Author     string `json:"author,omitempty"`
	Maintainer string `json:"maintainer,omitempty"`
	Cite       string `json:"cite,omitempty"`
	License    string `json:"license,omitempty"`

	Source string `json:"source"`
	Url    string `json:"url"`

	Custom map[string]interface{} `json:"custom,omitempty"`
}

// Enum for gear categories.
// Gear categories mainly affect how they are displayed in the user interfaces.
type GearCategory string

const (
	UtilityGear  GearCategory = "utility"
	AnalysisGear GearCategory = "analysis"

	// Legacy category; equivalent to Utility
	ConverterGear GearCategory = "converter"

	// Legacy category; equivalent to Utility
	QaGear GearCategory = "qa"
)

type GearDoc struct {
	Id       string       `json:"_id,omitempty"`
	Category GearCategory `json:"category,omitempty"`

	Gear   *Gear       `json:"gear,omitempty"`
	Source *GearSource `json:"exchange,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
}

func (c *Client) GetAllGears() ([]*GearDoc, *http.Response, error) {
	var aerr *Error
	var gears []*GearDoc
	resp, err := c.New().Get("gears").Receive(&gears, &aerr)
	return gears, resp, Coalesce(err, aerr)
}

func (c *Client) GetGear(id string) (*GearDoc, *http.Response, error) {
	var aerr *Error
	var gear *GearDoc
	resp, err := c.New().Get("gears/"+id).Receive(&gear, &aerr)
	return gear, resp, Coalesce(err, aerr)
}

func (c *Client) GetGearInvocation(id string) (map[string]interface{}, *http.Response, error) {
	var aerr *Error
	var response map[string]interface{}
	resp, err := c.New().Get("gears/"+id+"/invocation").Receive(&response, &aerr)
	return response, resp, Coalesce(err, aerr)
}

func (c *Client) AddGear(gear *GearDoc) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	gearName := gear.Gear.Name

	/*
		This language does not have type unions dot txt

		We need the json encoding annotation "omitempty" on every field, because otherwise modifications with unset fields might clobber. (Even if we don't currently support modifying gears.)

		But this becomes an issue when sending up a new gear: the API demands the config map be present, even if empty, and Go omits encoding a map with "omitempty" if the map is empty. So you can't send up a gear with no config, or with no inputs.

		So, remove "omitempty" off inputs & config, know that if we add modifications we'll need to check & strip those fields if empty, and explicitly set maps when adding if none were provided.
	*/

	if gear.Gear.Config == nil {
		gear.Gear.Config = map[string]map[string]interface{}{}
	}
	if gear.Gear.Inputs == nil {
		gear.Gear.Inputs = map[string]map[string]interface{}{}
	}

	resp, err := c.New().Post("gears/"+gearName).BodyJSON(gear).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

func (c *Client) DeleteGear(id string) (*http.Response, error) {
	var aerr *Error

	resp, err := c.New().Delete("gears/"+id).Receive(nil, &aerr)
	return resp, Coalesce(err, aerr)
}
