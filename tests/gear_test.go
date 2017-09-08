package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestGears() {
	gear := &api.Gear{
		Name:        RandStringLower(),
		Label:       RandString(),
		Description: RandString(),
		Version:     RandString(),
		Author:      RandString(),
		Maintainer:  RandString(),
		License:     "Other",
		Source:      "http://example.example",
		Url:         "http://example.example",
	}

	gearDoc := &api.GearDoc{
		Category: api.UtilityGear,
		Gear:     gear,
	}

	// Add
	gearId, _, err := t.AddGear(gearDoc)
	t.So(err, ShouldBeNil)

	// Get
	rGear, _, err := t.GetGear(gearId)
	t.So(err, ShouldBeNil)
	t.So(rGear.Gear.Name, ShouldEqual, gear.Name)
	now := time.Now()
	t.So(*rGear.Created, ShouldHappenBefore, now)
	t.So(*rGear.Modified, ShouldHappenBefore, now)

	// Get invocation
	gearSchema, _, err := t.GetGearInvocation(gearId)
	t.So(err, ShouldBeNil)
	t.So(gearSchema["$schema"], ShouldStartWith, "http://json-schema.org")

	// Get all
	gears, _, err := t.GetAllGears()
	t.So(err, ShouldBeNil)
	t.So(gears, ShouldContain, rGear)

	// Delete
	_, err = t.DeleteGear(gearId)
	t.So(err, ShouldBeNil)
	gears, _, err = t.GetAllGears()
	t.So(err, ShouldBeNil)
	t.So(gears, ShouldNotContain, rGear)
}

func (t *F) createTestGear() string {

	//
	// Do not modify the below gear document without checking the other callees!
	//

	gear := &api.Gear{
		Name:        RandStringLower(),
		Label:       RandString(),
		Description: RandString(),
		Version:     RandString(),
		Author:      RandString(),
		Maintainer:  RandString(),
		Inputs: map[string]map[string]interface{}{
			"any-file": {
				"base": "file",
			},
		},
		License: "Other",
		Source:  "http://example.example",
		Url:     "http://example.example",
	}

	gearDoc := &api.GearDoc{
		Category: "utility",
		Gear:     gear,
		Source: &api.GearSource{
			Commit:     "aex",
			RootfsHash: "sha384:oy",
			RootfsUrl:  "http://example.example",
		},
	}
	gearId, _, err := t.AddGear(gearDoc)
	t.So(err, ShouldBeNil)

	return gearId
}
