package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestAnalyses() {
	groupId, _, sessionId := t.createTestSession()
	gearId := t.createTestGear()

	src := UploadSourceFromString("yeats.txt", "A gaze blank and pitiless as the sun,")
	progress, resultChan := t.UploadToSession(sessionId, src)
	t.checkProgressChanEndsWith(progress, 37)
	t.So(<-resultChan, ShouldBeNil)

	analysis := &api.Analysis{
		Name:        RandString(),
		Description: RandString(),
	}

	filereference := &api.FileReference{
		Id:   sessionId,
		Type: "session",
		Name: "yeats.txt",
	}

	tag := RandString()

	job := &api.Job{
		GearId: gearId,
		Inputs: map[string]interface{}{
			"any-file": filereference,
		},
		Tags: []string{tag},
	}

	anaId, _, err := t.AddSessionAnalysis(sessionId, analysis, job)
	t.So(err, ShouldBeNil)

	session, _, err := t.GetSession(sessionId)
	t.So(err, ShouldBeNil)

	t.So(session.Analyses, ShouldHaveLength, 1)
	rAna := session.Analyses[0]

	t.So(rAna.Id, ShouldEqual, anaId)
	t.So(rAna.User, ShouldNotBeEmpty)
	t.So(rAna.Job.State, ShouldEqual, api.Pending)
	now := time.Now()
	t.So(*rAna.Created, ShouldHappenBefore, now)
	t.So(*rAna.Modified, ShouldHappenBefore, now)
	t.So(rAna.Files, ShouldHaveLength, 1)
	t.So(rAna.Files[0].Name, ShouldEqual, "yeats.txt")
	t.So(rAna.Files[0].Input, ShouldBeTrue)

	// Access analysis directly
	rAna2, _, err := t.GetAnalysis(rAna.Id)
	t.So(err, ShouldBeNil)
	t.So(rAna2, ShouldEqual, rAna2)

	// Access multiple analyses
	_, _, err = t.AddSessionAnalysis(sessionId, analysis, job)
	t.So(err, ShouldBeNil)

	// Try getting analysis incorrectly
	_, _, err = t.GetAnalyses("sessions", sessionId, "projects")
	t.So(err, ShouldNotBeNil)

	// Get all Session level analyses in group
	analyses, _, err := t.GetAnalyses("groups", groupId, "sessions")
	t.So(err, ShouldBeNil)
	t.So(len(analyses), ShouldEqual, 2)
	t.So(analyses[1], ShouldNotBeEmpty)

	// Get all Project level analyses in group (Will be zero)
	analyses, _, err = t.GetAnalyses("groups", groupId, "projects")
	t.So(err, ShouldBeNil)
	t.So(len(analyses), ShouldEqual, 0)

	// Run the job
	_, err = t.ChangeJobState(rAna.Job.Id, api.Running)
	t.So(err, ShouldBeNil)

	//
	// We can't test further than this, because /engine requires drone.
	//

	// Ad-hoc implementation of engine upload
	// src2 := UploadSourceFromString("yeats-result.txt", "And what rough beast, its hour come round at last,")
	// url := "engine?level=analysis&id=" + anaId + "&job=" + rAna.Job.Id
	// progress2, resultChan2 := t.UploadSimple(url, nil, src2)
	// t.checkProgressChanEndsWith(progress2, 50)
	// t.So(<-resultChan2, ShouldBeNil)

	// Check that the result file exists
	// session, _, err := t.GetSession(sessionId)
	// t.So(err, ShouldBeNil)

	// t.So(session.Analyses, ShouldHaveLength, 1)
	// rAna := session.Analyses[0]

	_, err = t.ChangeJobState(rAna.Job.Id, api.Complete)
	t.So(err, ShouldBeNil)

	// Analysis notes
	text := RandString()
	_, err = t.AddSessionAnalysisNote(sessionId, anaId, text)
	t.So(err, ShouldBeNil)

	// Check
	session, _, err = t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(session.Analyses, ShouldHaveLength, 1)
	rAna = session.Analyses[0]
	t.So(rAna.Notes, ShouldHaveLength, 1)
	t.So(rAna.Notes[0].UserId, ShouldNotBeEmpty)
	t.So(rAna.Notes[0].Text, ShouldEqual, text)
	now2 := time.Now()
	t.So(*rAna.Notes[0].Created, ShouldHappenBefore, now2)
	t.So(*rAna.Notes[0].Modified, ShouldHappenBefore, now2)
	t.So(*rAna.Modified, ShouldHappenAfter, now)
	t.So(*rAna.Modified, ShouldHappenBefore, now2)
}
