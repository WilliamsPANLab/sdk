package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestJobs() {
	_, _, _, acquisitionId := t.createTestAcquisition()
	gearId := t.createTestGear()

	tag := RandString()
	job := &api.Job{
		GearId: gearId,

		Destination: &api.ContainerReference{
			Id:   acquisitionId,
			Type: "acquisition",
		},

		Inputs: map[string]interface{}{
			"x": &api.FileReference{
				Id:   acquisitionId,
				Type: "acquisition",
				Name: "yeats.txt",
			},
		},

		Tags: []string{tag},
	}

	// Add
	jobId, _, err := t.AddJob(job)
	t.So(err, ShouldBeNil)

	// Get
	rJob, _, err := t.GetJob(jobId)
	t.So(err, ShouldBeNil)
	t.So(rJob.GearId, ShouldEqual, gearId)
	t.So(rJob.State, ShouldEqual, api.Pending)
	t.So(rJob.Attempt, ShouldEqual, 1)
	t.So(rJob.Origin, ShouldNotBeNil)
	t.So(rJob.Origin.Type, ShouldEqual, "user")
	t.So(rJob.Origin.Id, ShouldNotBeEmpty)
	t.So(rJob.Tags, ShouldContain, tag)
	now := time.Now()
	t.So(*rJob.Created, ShouldHappenBefore, now)
	t.So(*rJob.Modified, ShouldHappenBefore, now)

	// Modify
	tag2 := RandString()
	jobMod := &api.Job{
		Tags: []string{tag2},
	}

	// First as non-root, then as root
	_, err = t.ModifyJob(jobId, jobMod, false)
	t.So(err, ShouldNotBeNil)
	_, err = t.ModifyJob(jobId, jobMod, true)
	t.So(err, ShouldBeNil)

	// Check
	rJob, _, err = t.GetJob(jobId)
	t.So(err, ShouldBeNil)
	t.So(rJob.Tags, ShouldNotContain, tag)
	t.So(rJob.Tags, ShouldContain, tag2)
	t.So(rJob.State, ShouldEqual, api.Pending)

	// Cancel as non-root
	jobMod = &api.Job{State: api.Cancelled}
	_, err = t.ModifyJob(jobId, jobMod, false)
	t.So(err, ShouldBeNil)

	// Check
	rJob, _, err = t.GetJob(jobId)
	t.So(err, ShouldBeNil)
	t.So(rJob.State, ShouldEqual, api.Cancelled)
}

func (t *F) TestJobQueue() {
	_, _, _, acquisitionId := t.createTestAcquisition()
	gearId := t.createTestGear()

	tag := RandString()
	job := &api.Job{
		GearId: gearId,

		Destination: &api.ContainerReference{
			Id:   acquisitionId,
			Type: "acquisition",
		},

		Inputs: map[string]interface{}{
			"x": &api.FileReference{
				Id:   acquisitionId,
				Type: "acquisition",
				Name: "yeats.txt",
			},
		},

		Tags: []string{tag},
	}

	// Add
	jobId, _, err := t.AddJob(job)
	t.So(err, ShouldBeNil)

	// Check
	rJob, _, err := t.GetJob(jobId)
	t.So(err, ShouldBeNil)
	t.So(rJob.State, ShouldEqual, api.Pending)

	// Run
	jr, rJob, _, err := t.StartNextPendingJob(tag)
	t.So(err, ShouldBeNil)
	t.So(rJob, ShouldNotBeNil)
	t.So(jr, ShouldEqual, api.JobAquired)
	t.So(rJob.Id, ShouldEqual, jobId)
	t.So(rJob.Request, ShouldNotBeNil)
	t.So(rJob.Request.Target.Dir, ShouldEqual, "/flywheel/v0")

	// Next fetch with tag should not find any jobs
	jr, emptyJob, _, err := t.StartNextPendingJob(tag)
	t.So(err, ShouldBeNil)
	t.So(emptyJob, ShouldBeNil)
	t.So(jr, ShouldEqual, api.NoPendingJobs)

	// Heartbeat
	_, err = t.HeartbeatJob(jobId)
	t.So(err, ShouldBeNil)
	rJob2, _, err := t.GetJob(jobId)
	t.So(err, ShouldBeNil)
	t.So(*rJob2.Modified, ShouldHappenAfter, *rJob.Modified)

	// Finish
	_, err = t.ChangeJobState(jobId, api.Complete, true)
	t.So(err, ShouldBeNil)
	rJob3, _, err := t.GetJob(jobId)
	t.So(rJob3.State, ShouldEqual, api.Complete)
	t.So(*rJob3.Modified, ShouldHappenAfter, *rJob2.Modified)
}
