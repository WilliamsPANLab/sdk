package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

// Ref https://github.com/flywheel-io/sdk/issues/31
func (t *F) SkipTestSearch() {
	_, _, _, acquisitionId := t.createTestAcquisition()

	// Get Acquisition
	a, _, err := t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	// Ref https://github.com/flywheel-io/sdk/issues/32
	time.Sleep(1000 * time.Millisecond)

	s := &api.SearchQuery{
		ReturnType:   api.SessionString,
		SearchString: a.Name,
	}

	sR, _, err := t.Search(s)
	// t.So(a.Name, ShouldBeNil)
	t.So(err, ShouldBeNil)
	t.So(len(sR), ShouldEqual, 1)

	s = &api.SearchQuery{
		ReturnType:   api.AcquisitionString,
		SearchString: a.Name,
	}

	sR, _, err = t.Search(s)
	// t.So(a.Name, ShouldBeNil)
	t.So(err, ShouldBeNil)
	t.So(len(sR), ShouldEqual, 1)
}
