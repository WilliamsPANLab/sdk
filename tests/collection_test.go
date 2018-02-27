package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestCollections() {
	collectionName := RandString()

	collection := &api.Collection{
		Name:        collectionName,
		Description: RandString(),
	}

	// Add
	cId, _, err := t.AddCollection(collection)
	t.So(err, ShouldBeNil)

	// Get
	savedCollection, _, err := t.GetCollection(cId)
	t.So(err, ShouldBeNil)
	t.So(savedCollection.Id, ShouldEqual, cId)
	t.So(savedCollection.Name, ShouldEqual, collection.Name)
	now := time.Now()
	t.So(*savedCollection.Created, ShouldHappenBefore, now)
	t.So(*savedCollection.Modified, ShouldHappenBefore, now)

	// Add acquisition to the collection
	_, _, sessionId, acquisitionId := t.createTestAcquisition()
	_, err = t.AddAcquisitionsToCollection(cId, []string{acquisitionId})
	t.So(err, ShouldBeNil)

	// Get Sessions
	savedSessions, _, err := t.GetCollectionSessions(cId)
	t.So(savedSessions, ShouldHaveLength, 1)
	t.So(savedSessions[0].Id, ShouldEqual, sessionId)

	// Get Acquisitions
	savedAcquisitions, _, err := t.GetCollectionAcquisitions(cId)
	t.So(savedAcquisitions, ShouldHaveLength, 1)
	t.So(savedAcquisitions[0].Id, ShouldEqual, acquisitionId)

	// Get Session Acquisitions
	savedSessionAcquisitions, _, err := t.GetCollectionSessionAcquisitions(cId, savedSessions[0].Id)
	t.So(savedSessionAcquisitions, ShouldHaveLength, 1)
	t.So(savedSessionAcquisitions[0].Id, ShouldEqual, acquisitionId)

	// Add session to the collection
	_, _, sessionId, acquisitionId = t.createTestAcquisition()
	_, err = t.AddSessionsToCollection(cId, []string{sessionId})
	t.So(err, ShouldBeNil)

	// Get Sessions
	savedSessions, _, err = t.GetCollectionSessions(cId)
	t.So(savedSessions, ShouldHaveLength, 2)
	// Could add contains check

	// Get Acquisitions
	savedAcquisitions, _, err = t.GetCollectionAcquisitions(cId)
	t.So(savedAcquisitions, ShouldHaveLength, 2)
	// Could add contains check

	// Get Session Acquisitions
	savedSessionAcquisitions, _, err = t.GetCollectionSessionAcquisitions(cId, savedSessions[0].Id)
	t.So(savedSessionAcquisitions, ShouldHaveLength, 1)
	// Could add contains check

	// Get all
	collections, _, err := t.GetAllCollections()
	t.So(err, ShouldBeNil)
	// workaround: all-container endpoints skip some fields, single-container does not. this sets up the equality check
	savedCollection.Files = nil
	savedCollection.Notes = nil
	savedCollection.Info = nil
	// t.So(collections, ShouldContain, savedCollection)

	// Modify
	newName := RandString()
	collectionMod := &api.Collection{
		Name: newName,
	}
	_, err = t.ModifyCollection(cId, collectionMod)
	t.So(err, ShouldBeNil)

	// Check
	changedCollection, _, err := t.GetCollection(cId)
	t.So(changedCollection.Name, ShouldEqual, newName)
	t.So(*changedCollection.Created, ShouldBeSameTimeAs, *savedCollection.Created)
	t.So(*changedCollection.Modified, ShouldHappenAfter, *savedCollection.Modified)

	// Add note
	_, err = t.AddCollectionNote(cId, "This is a note")
	t.So(err, ShouldBeNil)
	changedCollection, _, err = t.GetCollection(cId)
	t.So(changedCollection.Notes, ShouldHaveLength, 1)
	t.So(changedCollection.Notes[0].Text, ShouldEqual, "This is a note")

	// Replace Info
	_, err = t.ReplaceCollectionInfo(cId, map[string]interface{}{
		"foo": 3,
		"bar": "qaz",
	})
	t.So(err, ShouldBeNil)

	// Set info
	_, err = t.SetCollectionInfo(cId, map[string]interface{}{
		"foo":   42,
		"hello": "world",
	})
	t.So(err, ShouldBeNil)

	changedCollection, _, err = t.GetCollection(cId)

	t.So(changedCollection.Info["foo"], ShouldEqual, 42)
	t.So(changedCollection.Info["bar"], ShouldEqual, "qaz")
	t.So(changedCollection.Info["hello"], ShouldEqual, "world")

	// Delete info fields
	_, err = t.DeleteCollectionInfoFields(cId, []string{"foo", "bar"})
	t.So(err, ShouldBeNil)

	changedCollection, _, err = t.GetCollection(cId)
	t.So(err, ShouldBeNil)

	t.So(changedCollection.Info["foo"], ShouldBeNil)
	t.So(changedCollection.Info["bar"], ShouldBeNil)
	t.So(changedCollection.Info["hello"], ShouldEqual, "world")

	// Delete
	_, err = t.DeleteCollection(cId)
	t.So(err, ShouldBeNil)
	collections, _, err = t.GetAllCollections()
	t.So(err, ShouldBeNil)
	t.So(collections, ShouldNotContain, savedCollection)
}

func (t *F) TestCollectionFiles() {
	collection := &api.Collection{Name: RandString()}
	collectionId, _, err := t.AddCollection(collection)
	t.So(err, ShouldBeNil)

	poem := "Things fall apart; the centre cannot hold;"
	t.uploadText(t.UploadToCollection, collectionId, "yeats.txt", poem)

	rCollection, _, err := t.GetCollection(collectionId)
	t.So(err, ShouldBeNil)
	t.So(rCollection.Files, ShouldHaveLength, 1)
	t.So(rCollection.Files[0].Name, ShouldEqual, "yeats.txt")
	t.So(rCollection.Files[0].Size, ShouldEqual, 42)
	t.So(rCollection.Files[0].Mimetype, ShouldEqual, "text/plain")

	// Download the same file and check content
	t.downloadText(t.DownloadFromCollection, collectionId, "yeats.txt", poem)

	// Test unauthorized download with ticket for the file
	t.downloadTextWithTicket(t.GetCollectionDownloadUrl, collectionId, "yeats.txt", poem)

	// Bundling: test file attributes
	t.So(rCollection.Files[0].Modality, ShouldEqual, "")
	t.So(rCollection.Files[0].Measurements, ShouldHaveLength, 0)
	t.So(rCollection.Files[0].Type, ShouldEqual, "text")

	_, response, err := t.ModifyCollectionFile(collectionId, "yeats.txt", &api.FileFields{
		Modality:     "MR",
		Measurements: []string{"functional"},
		Type:         "dicom",
	})
	t.So(err, ShouldBeNil)

	// Check that no jobs were triggered and attrs were modified
	t.So(response.JobsTriggered, ShouldEqual, 0)

	rCollection, _, err = t.GetCollection(collectionId)
	t.So(err, ShouldBeNil)
	t.So(rCollection.Files[0].Modality, ShouldEqual, "MR")
	t.So(rCollection.Files[0].Measurements, ShouldHaveLength, 1)
	t.So(rCollection.Files[0].Measurements[0], ShouldEqual, "functional")
	t.So(rCollection.Files[0].Type, ShouldEqual, "dicom")

	// Test file info
	t.So(rCollection.Files[0].Info, ShouldBeEmpty)
	_, err = t.ReplaceCollectionFileInfo(collectionId, "yeats.txt", map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	})
	t.So(err, ShouldBeNil)
	_, err = t.SetCollectionFileInfo(collectionId, "yeats.txt", map[string]interface{}{
		"c": 5,
	})
	t.So(err, ShouldBeNil)

	rCollection, _, err = t.GetCollection(collectionId)
	t.So(err, ShouldBeNil)
	t.So(rCollection.Files[0].Info["a"], ShouldEqual, 1)
	t.So(rCollection.Files[0].Info["b"], ShouldEqual, 2)
	t.So(rCollection.Files[0].Info["c"], ShouldEqual, 5)
	t.So(rCollection.Files[0].Info["d"], ShouldEqual, 4)

	_, err = t.DeleteCollectionFileInfoFields(collectionId, "yeats.txt", []string{"c", "d"})
	t.So(err, ShouldBeNil)

	rCollection, _, err = t.GetCollection(collectionId)
	t.So(err, ShouldBeNil)
	t.So(rCollection.Files[0].Info["a"], ShouldEqual, 1)
	t.So(rCollection.Files[0].Info["b"], ShouldEqual, 2)
	t.So(rCollection.Files[0].Info["c"], ShouldBeNil)
	t.So(rCollection.Files[0].Info["d"], ShouldBeNil)

	_, err = t.ReplaceCollectionFileInfo(collectionId, "yeats.txt", map[string]interface{}{})
	rCollection, _, err = t.GetCollection(collectionId)
	t.So(err, ShouldBeNil)
	t.So(rCollection.Files[0].Info, ShouldBeEmpty)

	// Delete file
	_, err = t.DeleteCollectionFile(collectionId, "yeats.txt")
	t.So(err, ShouldBeNil)

	rCollection, _, err = t.GetCollection(collectionId)
	t.So(err, ShouldBeNil)
	t.So(len(rCollection.Files), ShouldEqual, 0)
}

func (t *F) TestCollectionAnalysis() {
	collection := &api.Collection{Name: RandString()}
	collectionId, _, err := t.AddCollection(collection)
	t.So(err, ShouldBeNil)

	poem := "Things fall apart; the centre cannot hold;"
	t.uploadText(t.UploadToCollection, collectionId, "yeats.txt", poem)

	filereference := &api.FileReference{
		Id:   collectionId,
		Type: "collection",
		Name: "yeats.txt",
	}

	analysis := &api.AdhocAnalysis{
		Name:        RandString(),
		Description: RandString(),
		Inputs:      []*api.FileReference{filereference},
	}

	anaId, _, err := t.AddCollectionAnalysis(collectionId, analysis)
	t.So(err, ShouldBeNil)
	t.So(anaId, ShouldNotBeNil)

	rCollection, _, err := t.GetCollection(collectionId)
	t.So(err, ShouldBeNil)
	t.So(rCollection.Id, ShouldEqual, collectionId)

	t.So(rCollection.Analyses, ShouldHaveLength, 1)
	rAna := rCollection.Analyses[0]

	t.So(rAna.Id, ShouldEqual, anaId)
	t.So(rAna.User, ShouldNotBeEmpty)
	t.So(rAna.Job, ShouldBeNil)
	now := time.Now()
	t.So(*rAna.Created, ShouldHappenBefore, now)
	t.So(*rAna.Modified, ShouldHappenBefore, now)
	t.So(rAna.Inputs, ShouldHaveLength, 1)
	t.So(rAna.Inputs[0].Name, ShouldEqual, "yeats.txt")
}
