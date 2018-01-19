package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestAcquisitions() {
	_, _, sessionId := t.createTestSession()

	acquisitionName := RandString()
	acquisition := &api.Acquisition{
		Name:      acquisitionName,
		SessionId: sessionId,
	}

	// Add
	acquisitionId, _, err := t.AddAcquisition(acquisition)
	t.So(err, ShouldBeNil)

	// Get
	rAcquisition, _, err := t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	t.So(rAcquisition.Id, ShouldEqual, acquisitionId)
	t.So(rAcquisition.Name, ShouldEqual, acquisition.Name)
	now := time.Now()
	t.So(*rAcquisition.Created, ShouldHappenBefore, now)
	t.So(*rAcquisition.Modified, ShouldHappenBefore, now)

	// Get all
	acquisitions, _, err := t.GetAllAcquisitions()
	t.So(err, ShouldBeNil)
	// workaround: all-container endpoints skip some fields, single-container does not. this sets up the equality check
	rAcquisition.Files = []*api.File{}
	rAcquisition.Notes = []*api.Note{}
	rAcquisition.Tags = []string{}
	rAcquisition.Info = map[string]interface{}{}
	t.So(acquisitions, ShouldContain, rAcquisition)

	// Get from parent
	acquisitions, _, err = t.GetSessionAcquisitions(sessionId)
	t.So(err, ShouldBeNil)
	t.So(acquisitions, ShouldContain, rAcquisition)

	// Modify
	newName := RandString()
	acquisitionMod := &api.Acquisition{
		Name: newName,
	}
	_, err = t.ModifyAcquisition(acquisitionId, acquisitionMod)
	t.So(err, ShouldBeNil)
	changedAcquisition, _, err := t.GetAcquisition(acquisitionId)
	t.So(changedAcquisition.Name, ShouldEqual, newName)
	t.So(*changedAcquisition.Created, ShouldBeSameTimeAs, *rAcquisition.Created)
	t.So(*changedAcquisition.Modified, ShouldHappenAfter, *rAcquisition.Modified)

	// Notes, tags
	message := "This is a note"
	_, err = t.AddAcquisitionNote(acquisitionId, message)
	t.So(err, ShouldBeNil)
	tag := "example-tag"
	_, err = t.AddAcquisitionTag(acquisitionId, tag)
	t.So(err, ShouldBeNil)

	// Replace Info
	_, err = t.ReplaceAcquisitionInfo(acquisitionId, map[string]interface{}{
		"foo": 3,
		"bar": "qaz",
	})
	t.So(err, ShouldBeNil)

	// Set info
	_, err = t.SetAcquisitionInfo(acquisitionId, map[string]interface{}{
		"foo":   42,
		"hello": "world",
	})
	t.So(err, ShouldBeNil)

	// Check
	rAcquisition, _, err = t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	t.So(rAcquisition.Notes, ShouldHaveLength, 1)
	t.So(rAcquisition.Notes[0].Text, ShouldEqual, message)
	t.So(rAcquisition.Tags, ShouldHaveLength, 1)
	t.So(rAcquisition.Tags[0], ShouldEqual, tag)

	t.So(rAcquisition.Info["foo"], ShouldEqual, 42)
	t.So(rAcquisition.Info["bar"], ShouldEqual, "qaz")
	t.So(rAcquisition.Info["hello"], ShouldEqual, "world")

	// Delete info fields
	_, err = t.DeleteAcquisitionInfoFields(acquisitionId, []string{"foo", "bar"})
	t.So(err, ShouldBeNil)

	rAcquisition, _, err = t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)

	t.So(rAcquisition.Info["foo"], ShouldBeNil)
	t.So(rAcquisition.Info["bar"], ShouldBeNil)
	t.So(rAcquisition.Info["hello"], ShouldEqual, "world")

	// Delete
	_, err = t.DeleteAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	acquisitions, _, err = t.GetAllAcquisitions()
	t.So(err, ShouldBeNil)
	t.So(acquisitions, ShouldNotContain, rAcquisition)
}

func (t *F) TestAcquisitionFiles() {
	_, _, sessionId := t.createTestSession()

	acquisition := &api.Acquisition{Name: RandString(), SessionId: sessionId}
	acquisitionId, _, err := t.AddAcquisition(acquisition)
	t.So(err, ShouldBeNil)

	poem := "Turning and turning in the widening gyre"
	t.uploadText(t.UploadToAcquisition, acquisitionId, "yeats.txt", poem)

	rAcquisition, _, err := t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	t.So(rAcquisition.Files, ShouldHaveLength, 1)
	t.So(rAcquisition.Files[0].Name, ShouldEqual, "yeats.txt")
	t.So(rAcquisition.Files[0].Size, ShouldEqual, 40)
	t.So(rAcquisition.Files[0].Mimetype, ShouldEqual, "text/plain")

	// Download the same file and check content
	t.downloadText(t.DownloadFromAcquisition, acquisitionId, "yeats.txt", poem)

	// Test unauthorized download with ticket for the file
	t.downloadTextWithTicket(t.GetAcquisitionDownloadUrl, acquisitionId, "yeats.txt", poem)

	// Bundling: test file attributes
	t.So(rAcquisition.Files[0].Modality, ShouldEqual, "")
	t.So(rAcquisition.Files[0].Measurements, ShouldHaveLength, 0)
	t.So(rAcquisition.Files[0].Type, ShouldEqual, "text")

	_, response, err := t.ModifyAcquisitionFile(acquisitionId, "yeats.txt", &api.FileFields{
		Modality:     "MR",
		Measurements: []string{"functional"},
		Type:         "dicom",
	})
	t.So(err, ShouldBeNil)

	// Check that no jobs were triggered and attrs were modified
	t.So(response.JobsTriggered, ShouldEqual, 0)

	rAcquisition, _, err = t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	t.So(rAcquisition.Files[0].Modality, ShouldEqual, "MR")
	t.So(rAcquisition.Files[0].Measurements, ShouldHaveLength, 1)
	t.So(rAcquisition.Files[0].Measurements[0], ShouldEqual, "functional")
	t.So(rAcquisition.Files[0].Type, ShouldEqual, "dicom")

	// Test file info
	t.So(rAcquisition.Files[0].Info, ShouldBeEmpty)
	_, err = t.ReplaceAcquisitionFileInfo(acquisitionId, "yeats.txt", map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	})
	t.So(err, ShouldBeNil)
	_, err = t.SetAcquisitionFileInfo(acquisitionId, "yeats.txt", map[string]interface{}{
		"c": 5,
	})
	t.So(err, ShouldBeNil)

	rAcquisition, _, err = t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	t.So(rAcquisition.Files[0].Info["a"], ShouldEqual, 1)
	t.So(rAcquisition.Files[0].Info["b"], ShouldEqual, 2)
	t.So(rAcquisition.Files[0].Info["c"], ShouldEqual, 5)
	t.So(rAcquisition.Files[0].Info["d"], ShouldEqual, 4)

	_, err = t.DeleteAcquisitionFileInfoFields(acquisitionId, "yeats.txt", []string{"c", "d"})
	t.So(err, ShouldBeNil)

	rAcquisition, _, err = t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	t.So(rAcquisition.Files[0].Info["a"], ShouldEqual, 1)
	t.So(rAcquisition.Files[0].Info["b"], ShouldEqual, 2)
	t.So(rAcquisition.Files[0].Info["c"], ShouldBeNil)
	t.So(rAcquisition.Files[0].Info["d"], ShouldBeNil)

	_, err = t.ReplaceAcquisitionFileInfo(acquisitionId, "yeats.txt", map[string]interface{}{})
	rAcquisition, _, err = t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	t.So(rAcquisition.Files[0].Info, ShouldBeEmpty)
}

func (t *F) createTestAcquisition() (string, string, string, string) {
	groupId, projectId, sessionId := t.createTestSession()

	acquisitionName := RandString()
	acquisition := &api.Acquisition{
		Name:      acquisitionName,
		SessionId: sessionId,
	}
	acquisitionId, _, err := t.AddAcquisition(acquisition)
	t.So(err, ShouldBeNil)

	return groupId, projectId, sessionId, acquisitionId
}
