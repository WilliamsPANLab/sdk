package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestBatch() {
	_, _, _, acquisitionId := t.createTestAcquisition()
	gearId := t.createTestGear()

	poem := "The falcon cannot hear the falconer;"
	t.uploadText(t.UploadToAcquisition, acquisitionId, "yeats.txt", poem)

	// Add
	tag := RandString()
	targets := []*api.ContainerReference{
		{
			Id:   acquisitionId,
			Type: "acquisition",
		},
	}
	proposal, _, err := t.ProposeBatch(gearId, nil, []string{tag}, targets)
	t.So(err, ShouldBeNil)
	t.So(proposal.Id, ShouldNotBeEmpty)
	t.So(proposal.GearId, ShouldEqual, gearId)
	t.So(proposal.Origin, ShouldNotBeNil)
	t.So(proposal.Origin.Type, ShouldEqual, "user")
	t.So(proposal.Origin.Id, ShouldNotBeEmpty)
	t.So(proposal.Matched, ShouldHaveLength, 1)
	t.So(proposal.Ambiguous, ShouldBeEmpty)
	t.So(proposal.MissingPermissions, ShouldBeEmpty)
	t.So(proposal.NotMatched, ShouldBeEmpty)
	now := time.Now()
	t.So(*proposal.Created, ShouldHappenBefore, now)
	t.So(*proposal.Modified, ShouldHappenBefore, now)

	// Get
	rBatch, _, err := t.GetBatch(proposal.Id)
	t.So(err, ShouldBeNil)
	t.So(rBatch.GearId, ShouldEqual, proposal.GearId)
	t.So(rBatch.State, ShouldEqual, api.Pending)

	// Get all
	batches, _, err := t.GetAllBatches()
	t.So(err, ShouldBeNil)
	t.So(batches, ShouldContain, rBatch)

	JobQueue.Lock()
	defer JobQueue.Unlock()

	// Start
	jobs, _, err := t.StartBatch(proposal.Id)
	t.So(err, ShouldBeNil)
	t.So(jobs, ShouldHaveLength, 1)

	// Get again
	rBatch2, _, err := t.GetBatch(proposal.Id)
	t.So(err, ShouldBeNil)
	t.So(rBatch2.State, ShouldEqual, api.Running)
	t.So(*rBatch2.Modified, ShouldHappenAfter, *rBatch.Modified)

	// Cancel
	cancelled, _, err := t.CancelBatch(proposal.Id)
	t.So(err, ShouldBeNil)
	t.So(cancelled, ShouldEqual, 1)
}
