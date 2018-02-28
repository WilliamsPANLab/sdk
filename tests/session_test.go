package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestSessions() {
	_, projectId := t.createTestProject()

	sessionName := RandString()
	session := &api.Session{
		Name:      sessionName,
		ProjectId: projectId,
		Info: map[string]interface{}{
			"some-key": 37,
		},
		Subject: &api.Subject{
			Code:      RandStringLower(),
			Firstname: RandString(),
			Lastname:  RandString(),
			Sex:       "other",
			Age:       56,
			Info: map[string]interface{}{
				"some-subject-key": 37,
			},
		},
	}

	// Add
	sessionId, _, err := t.AddSession(session)
	t.So(err, ShouldBeNil)

	// Get
	rSession, _, err := t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(rSession.Id, ShouldEqual, sessionId)
	t.So(rSession.Name, ShouldEqual, session.Name)
	now := time.Now()
	t.So(rSession.Info, ShouldContainKey, "some-key")
	t.So(rSession.Info["some-key"], ShouldEqual, 37)
	t.So(*rSession.Created, ShouldHappenBefore, now)
	t.So(*rSession.Modified, ShouldHappenBefore, now)
	t.So(*rSession.Subject, ShouldNotBeNil)
	t.So(rSession.Subject.Id, ShouldNotBeEmpty)
	t.So(rSession.Subject.Firstname, ShouldResemble, session.Subject.Firstname)

	// Get all
	sessions, _, err := t.GetAllSessions()
	t.So(err, ShouldBeNil)
	// workaround: all-container endpoints skip some fields, single-container does not. this sets up the equality check
	rSession.Files = []*api.File{}
	rSession.Notes = []*api.Note{}
	rSession.Tags = []string{}
	rSession.Info = map[string]interface{}{}
	rSession.Analyses = nil
	rSession.Subject = &api.Subject{
		Id:   rSession.Subject.Id,
		Code: rSession.Subject.Code,
		Info: map[string]interface{}{},
	}
	t.So(sessions, ShouldContain, rSession)

	// Get from parent
	sessions, _, err = t.GetProjectSessions(projectId)
	t.So(err, ShouldBeNil)
	t.So(sessions, ShouldContain, rSession)

	// Modify
	newName := RandString()
	sessionMod := &api.Session{
		Name: newName,
	}
	_, err = t.ModifySession(sessionId, sessionMod)
	t.So(err, ShouldBeNil)
	changedSession, _, err := t.GetSession(sessionId)
	t.So(changedSession.Name, ShouldEqual, newName)
	t.So(*changedSession.Created, ShouldBeSameTimeAs, *rSession.Created)
	t.So(*changedSession.Modified, ShouldHappenAfter, *rSession.Modified)

	// Notes, tags
	message := "This is a note"
	_, err = t.AddSessionNote(sessionId, message)
	t.So(err, ShouldBeNil)
	tag := "example-tag"
	_, err = t.AddSessionTag(sessionId, tag)
	t.So(err, ShouldBeNil)

	// Replace Info
	_, err = t.ReplaceSessionInfo(sessionId, map[string]interface{}{
		"foo": 3,
		"bar": "qaz",
	})
	t.So(err, ShouldBeNil)

	// Set info
	_, err = t.SetSessionInfo(sessionId, map[string]interface{}{
		"foo":   42,
		"hello": "world",
	})
	t.So(err, ShouldBeNil)

	// Check
	rSession, _, err = t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(rSession.Notes, ShouldHaveLength, 1)
	t.So(rSession.Notes[0].Text, ShouldEqual, message)
	t.So(rSession.Tags, ShouldHaveLength, 1)
	t.So(rSession.Tags[0], ShouldEqual, tag)

	t.So(rSession.Info["foo"], ShouldEqual, 42)
	t.So(rSession.Info["bar"], ShouldEqual, "qaz")
	t.So(rSession.Info["hello"], ShouldEqual, "world")

	// Delete info fields
	_, err = t.DeleteSessionInfoFields(sessionId, []string{"foo", "bar"})
	t.So(err, ShouldBeNil)

	rSession, _, err = t.GetSession(sessionId)
	t.So(err, ShouldBeNil)

	t.So(rSession.Info["foo"], ShouldBeNil)
	t.So(rSession.Info["bar"], ShouldBeNil)
	t.So(rSession.Info["hello"], ShouldEqual, "world")

	// Delete
	_, err = t.DeleteSession(sessionId)
	t.So(err, ShouldBeNil)
	sessions, _, err = t.GetAllSessions()
	t.So(err, ShouldBeNil)
	t.So(sessions, ShouldNotContain, rSession)
}

func (t *F) TestSessionFiles() {
	_, projectId := t.createTestProject()
	session := &api.Session{Name: RandString(), ProjectId: projectId}
	sessionId, _, err := t.AddSession(session)
	t.So(err, ShouldBeNil)

	poem := "The best lack all conviction, while the worst"
	t.uploadText(t.UploadToSession, sessionId, "yeats.txt", poem)

	rSession, _, err := t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(rSession.Files, ShouldHaveLength, 1)
	t.So(rSession.Files[0].Name, ShouldEqual, "yeats.txt")
	t.So(rSession.Files[0].Size, ShouldEqual, 45)
	t.So(rSession.Files[0].Mimetype, ShouldEqual, "text/plain")

	// Download the same file and check content
	t.downloadText(t.DownloadFromSession, sessionId, "yeats.txt", poem)

	// Test unauthorized download with ticket for the file
	t.downloadTextWithTicket(t.GetSessionDownloadUrl, sessionId, "yeats.txt", poem)

	// Bundling: test file attributes
	t.So(rSession.Files[0].Modality, ShouldEqual, "")
	t.So(rSession.Files[0].Measurements, ShouldHaveLength, 0)
	t.So(rSession.Files[0].Type, ShouldEqual, "text")

	_, response, err := t.ModifySessionFile(sessionId, "yeats.txt", &api.FileFields{
		Modality:     "modality",
		Measurements: []string{"measurement"},
		Type:         "type",
	})
	t.So(err, ShouldBeNil)

	// Check that no jobs were triggered and attrs were modified
	t.So(response.JobsTriggered, ShouldEqual, 0)

	rSession, _, err = t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(rSession.Files[0].Modality, ShouldEqual, "modality")
	t.So(rSession.Files[0].Measurements, ShouldHaveLength, 1)
	t.So(rSession.Files[0].Measurements[0], ShouldEqual, "measurement")
	t.So(rSession.Files[0].Type, ShouldEqual, "type")

	// Test file info
	t.So(rSession.Files[0].Info, ShouldBeEmpty)
	_, err = t.ReplaceSessionFileInfo(sessionId, "yeats.txt", map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	})
	t.So(err, ShouldBeNil)
	_, err = t.SetSessionFileInfo(sessionId, "yeats.txt", map[string]interface{}{
		"c": 5,
	})
	t.So(err, ShouldBeNil)

	rSession, _, err = t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(rSession.Files[0].Info["a"], ShouldEqual, 1)
	t.So(rSession.Files[0].Info["b"], ShouldEqual, 2)
	t.So(rSession.Files[0].Info["c"], ShouldEqual, 5)
	t.So(rSession.Files[0].Info["d"], ShouldEqual, 4)

	_, err = t.DeleteSessionFileInfoFields(sessionId, "yeats.txt", []string{"c", "d"})
	t.So(err, ShouldBeNil)

	rSession, _, err = t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(rSession.Files[0].Info["a"], ShouldEqual, 1)
	t.So(rSession.Files[0].Info["b"], ShouldEqual, 2)
	t.So(rSession.Files[0].Info["c"], ShouldBeNil)
	t.So(rSession.Files[0].Info["d"], ShouldBeNil)

	_, err = t.ReplaceSessionFileInfo(sessionId, "yeats.txt", map[string]interface{}{})
	rSession, _, err = t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(rSession.Files[0].Info, ShouldBeEmpty)

	// Delete file
	_, err = t.DeleteSessionFile(sessionId, "yeats.txt")
	t.So(err, ShouldBeNil)

	rSession, _, err = t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(len(rSession.Files), ShouldEqual, 0)
}

func (t *F) createTestSession() (string, string, string) {
	groupId, projectId := t.createTestProject()

	sessionName := RandString()
	session := &api.Session{
		Name:      sessionName,
		ProjectId: projectId,
		Info: map[string]interface{}{
			"some-key": 37,
		},
		Subject: &api.Subject{
			Code:      RandStringLower(),
			Firstname: RandString(),
			Lastname:  RandString(),
			Sex:       "other",
			Age:       56,
			Info: map[string]interface{}{
				"some-subject-key": 37,
			},
		},
	}
	sessionId, _, err := t.AddSession(session)
	t.So(err, ShouldBeNil)

	return groupId, projectId, sessionId
}
