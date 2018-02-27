package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestAnalysis() {
	groupId, _, sessionId := t.createTestSession()

	src := UploadSourceFromString("yeats.txt", "A gaze blank and pitiless as the sun,")
	progress, resultChan := t.UploadToSession(sessionId, src)
	t.checkProgressChanEndsWith(progress, 37)
	t.So(<-resultChan, ShouldBeNil)

	filereference := &api.FileReference{
		Id:   sessionId,
		Type: "session",
		Name: "yeats.txt",
	}

	analysis := &api.AdhocAnalysis{
		Name:        RandString(),
		Description: RandString(),
		Inputs:      []*api.FileReference{filereference},
	}

	anaId, _, err := t.AddSessionAnalysis(sessionId, analysis)
	t.So(err, ShouldBeNil)

	session, _, err := t.GetSession(sessionId)
	t.So(err, ShouldBeNil)

	t.So(session.Analyses, ShouldHaveLength, 1)
	rAna := session.Analyses[0]

	t.So(rAna.Id, ShouldEqual, anaId)
	t.So(rAna.User, ShouldNotBeEmpty)
	t.So(rAna.Job, ShouldBeNil)
	now := time.Now()
	t.So(*rAna.Created, ShouldHappenBefore, now)
	t.So(*rAna.Modified, ShouldHappenBefore, now)
	t.So(rAna.Inputs, ShouldHaveLength, 1)
	t.So(rAna.Inputs[0].Name, ShouldEqual, "yeats.txt")

	// Access analysis directly
	rAna2, _, err := t.GetAnalysis(rAna.Id)
	t.So(err, ShouldBeNil)
	t.So(rAna2, ShouldEqual, rAna2)

	// Analysis notes
	text := RandString()
	_, err = t.AddAnalysisNote(anaId, text)
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

	// Access multiple analyses
	_, _, err = t.AddSessionAnalysis(sessionId, analysis)
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

	// Notes, tags
	tag := "example-tag"
	_, err = t.AddAnalysisTag(anaId, tag)
	t.So(err, ShouldBeNil)

	// Replace Info
	_, err = t.ReplaceAnalysisInfo(anaId, map[string]interface{}{
		"foo": 3,
		"bar": "qaz",
	})
	t.So(err, ShouldBeNil)

	// Set info
	_, err = t.SetAnalysisInfo(anaId, map[string]interface{}{
		"foo":   42,
		"hello": "world",
	})
	t.So(err, ShouldBeNil)

	// Check
	rAna, _, err = t.GetAnalysis(anaId)
	t.So(err, ShouldBeNil)
	t.So(rAna.Tags, ShouldHaveLength, 1)
	t.So(rAna.Tags[0], ShouldEqual, tag)

	t.So(rAna.Info["foo"], ShouldEqual, 42)
	t.So(rAna.Info["bar"], ShouldEqual, "qaz")
	t.So(rAna.Info["hello"], ShouldEqual, "world")

	// Delete info fields
	_, err = t.DeleteAnalysisInfoFields(anaId, []string{"foo", "bar"})
	t.So(err, ShouldBeNil)

	rAna, _, err = t.GetAnalysis(anaId)
	t.So(err, ShouldBeNil)

	t.So(rAna.Info["foo"], ShouldBeNil)
	t.So(rAna.Info["bar"], ShouldBeNil)
	t.So(rAna.Info["hello"], ShouldEqual, "world")
}

func (t *F) TestAnalysisFiles() {
	_, _, sessionId := t.createTestSession()

	poemIn := "A gaze blank and pitiless as the sun,"

	src := UploadSourceFromString("yeats.txt", poemIn)
	progress, resultChan := t.UploadToSession(sessionId, src)
	t.checkProgressChanEndsWith(progress, 37)
	t.So(<-resultChan, ShouldBeNil)

	filereference := &api.FileReference{
		Id:   sessionId,
		Type: "session",
		Name: "yeats.txt",
	}

	analysis := &api.AdhocAnalysis{
		Name:        RandString(),
		Description: RandString(),
		Inputs:      []*api.FileReference{filereference},
	}

	analysisId, _, err := t.AddSessionAnalysis(sessionId, analysis)

	// Download the input file and check content
	t.downloadText(t.DownloadInputFromAnalysis, analysisId, "yeats.txt", poemIn)

	// Test unauthorized download with ticket for the file
	t.downloadTextWithTicket(t.GetAnalysisInputDownloadUrl, analysisId, "yeats.txt", poemIn)

	poemOut := "Surely the Second Coming is at hand."
	t.uploadText(t.UploadToAnalysis, analysisId, "yeats-out.txt", poemOut)

	rAnalysis, _, err := t.GetAnalysis(analysisId)
	t.So(err, ShouldBeNil)
	t.So(rAnalysis.Files, ShouldHaveLength, 1)
	t.So(rAnalysis.Files[0].Name, ShouldEqual, "yeats-out.txt")
	t.So(rAnalysis.Files[0].Size, ShouldEqual, 36)
	t.So(rAnalysis.Files[0].Mimetype, ShouldEqual, "text/plain")

	// Download the same file and check content
	t.downloadText(t.DownloadFromAnalysis, analysisId, "yeats-out.txt", poemOut)

	// Test unauthorized download with ticket for the file
	t.downloadTextWithTicket(t.GetAnalysisDownloadUrl, analysisId, "yeats-out.txt", poemOut)

	// File metadata / info modification not supported for analyses
}
