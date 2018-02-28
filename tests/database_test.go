package tests

import (
	"encoding/hex"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestDbAccess() {
	// t.SkipNow is not provided; for now just silently succeed.
	if t.MongoString == "" {
		return
	}

	// Connect
	m, err := mgo.DialWithTimeout(t.MongoString, time.Duration(3*time.Second))
	t.So(err, ShouldBeNil)
	defer m.Close()
	m.SetSafe(&mgo.Safe{}) //durability

	// Hierarchy
	groupId, projectId, sessionId, acquisitionId := t.createTestAcquisition()
	gearId := t.createTestGear()

	//
	// Check deserialization
	//

	var group api.Group
	err = m.DB("scitran").C("groups").Find(t.createStringIdQuery(groupId)).One(&group)
	t.So(err, ShouldBeNil)
	t.So(group.Id, ShouldEqual, groupId)

	var project api.Project
	err = m.DB("scitran").C("projects").Find(t.createIdQuery(projectId)).One(&project)
	t.So(err, ShouldBeNil)
	t.checkBinaryId(project.Id, projectId)

	var session api.Session
	err = m.DB("scitran").C("sessions").Find(t.createIdQuery(sessionId)).One(&session)
	t.So(err, ShouldBeNil)
	t.checkBinaryId(session.Id, sessionId)

	var acquisition api.Acquisition
	err = m.DB("scitran").C("acquisitions").Find(t.createIdQuery(acquisitionId)).One(&acquisition)
	t.So(err, ShouldBeNil)
	t.checkBinaryId(acquisition.Id, acquisitionId)

	rUser, _, err := t.GetCurrentUser()
	t.So(err, ShouldBeNil)

	var user api.User
	err = m.DB("scitran").C("users").Find(t.createStringIdQuery(rUser.Id)).One(&user)
	t.So(err, ShouldBeNil)

	//
	// Compare to same values in API
	//

	savedGroup, _, err := t.GetGroup(groupId)
	t.So(err, ShouldBeNil)
	t.So(savedGroup.Name, ShouldEqual, group.Name)
	t.So(*savedGroup.Created, ShouldBeSameTimeAs, *group.Created)
	t.So(*savedGroup.Modified, ShouldBeSameTimeAs, *group.Modified)

	rProject, _, err := t.GetProject(projectId)
	t.So(err, ShouldBeNil)
	t.So(rProject.Id, ShouldEqual, projectId)
	t.So(rProject.Name, ShouldEqual, project.Name)
	t.So(rProject.Description, ShouldEqual, project.Description)
	t.So(rProject.Info, ShouldContainKey, "some-key")
	t.So(rProject.Info["some-key"], ShouldEqual, 37)
	t.So(*rProject.Created, ShouldBeSameTimeAs, *project.Created)
	t.So(*rProject.Modified, ShouldBeSameTimeAs, *project.Modified)

	rSession, _, err := t.GetSession(sessionId)
	t.So(err, ShouldBeNil)
	t.So(rSession.Id, ShouldEqual, sessionId)
	t.So(rSession.Name, ShouldEqual, session.Name)
	t.So(rSession.Info, ShouldContainKey, "some-key")
	t.So(rSession.Info["some-key"], ShouldEqual, 37)
	t.So(*rSession.Created, ShouldBeSameTimeAs, *session.Created)
	t.So(*rSession.Modified, ShouldBeSameTimeAs, *session.Modified)
	t.So(*rSession.Subject, ShouldNotBeNil)
	t.So(rSession.Subject.Id, ShouldNotBeEmpty)
	t.So(rSession.Subject.Firstname, ShouldResemble, session.Subject.Firstname)

	rAcquisition, _, err := t.GetAcquisition(acquisitionId)
	t.So(err, ShouldBeNil)
	t.So(rAcquisition.Id, ShouldEqual, acquisitionId)
	t.So(rAcquisition.Name, ShouldEqual, acquisition.Name)
	t.So(*rAcquisition.Created, ShouldBeSameTimeAs, *acquisition.Created)
	t.So(*rAcquisition.Modified, ShouldBeSameTimeAs, *acquisition.Modified)

	t.So(user.Id, ShouldEqual, rUser.Id)
	t.So(user.Email, ShouldEqual, rUser.Email)
	t.So(user.Firstname, ShouldEqual, rUser.Firstname)
	t.So(user.Lastname, ShouldEqual, rUser.Lastname)
	t.So(*user.Created, ShouldBeSameTimeAs, *rUser.Created)
	t.So(*user.Modified, ShouldBeSameTimeAs, *rUser.Modified)

	//
	// Check jobs
	//

	src := UploadSourceFromString("yeats.txt", "A gaze blank and pitiless as the sun,")
	progress, resultChan := t.UploadToSession(sessionId, src)
	t.checkProgressChanEndsWith(progress, 37)
	t.So(<-resultChan, ShouldBeNil)

	filereference := &api.FileReference{
		Id:   sessionId,
		Type: "session",
		Name: "yeats.txt",
	}

	tag := RandString()
	sendJob := &api.Job{
		GearId: gearId,
		Inputs: map[string]interface{}{
			"any-file": filereference,
		},
		Tags: []string{tag},
	}
	jobId, _, err := t.AddJob(sendJob)
	t.So(err, ShouldBeNil)

	var job api.Job
	err = m.DB("scitran").C("jobs").Find(t.createIdQuery(jobId)).One(&job)
	t.So(err, ShouldBeNil)
	t.checkBinaryId(job.Id, jobId)

	rJob, _, err := t.GetJob(jobId)
	t.So(err, ShouldBeNil)
	t.So(rJob.GearId, ShouldEqual, gearId)
	t.So(rJob.State, ShouldEqual, api.Pending)
	t.So(rJob.Attempt, ShouldEqual, 1)
	t.So(rJob.Origin, ShouldNotBeNil)
	t.So(rJob.Origin.Type, ShouldEqual, "user")
	t.So(rJob.Origin.Id, ShouldNotBeEmpty)
	t.So(rJob.Tags, ShouldContain, tag)
	t.So(*rJob.Created, ShouldBeSameTimeAs, *job.Created)
	t.So(*rJob.Modified, ShouldBeSameTimeAs, *job.Modified)

	var gear api.GearDoc
	err = m.DB("scitran").C("gears").Find(t.createIdQuery(gearId)).One(&gear)
	t.So(err, ShouldBeNil)
	t.checkBinaryId(gear.Id, gearId)

	rGear, _, err := t.GetGear(gearId)
	t.So(err, ShouldBeNil)
	t.So(rGear.Gear.Name, ShouldEqual, gear.Gear.Name)
	t.So(*rGear.Created, ShouldBeSameTimeAs, *gear.Created)
	t.So(*rGear.Modified, ShouldBeSameTimeAs, *gear.Modified)

}

// Convert a dumb binary string into hex, and check that it's DB-ID-ish.
func (t *F) checkBinaryId(binaryString, compare string) {
	// If the string is empty, we're missing a bson annotation
	t.So(binaryString, ShouldNotBeEmpty)

	// Initial data should be binary
	t.So(binaryString, ShouldNotMatchRegex, idRegex)

	// Encode the garbage into something useful
	id := hex.EncodeToString([]byte(binaryString))

	// Now it should be DB-ID-ish.
	t.So(id, ShouldMatchRegex, idRegex)

	// Result should be as expected
	t.So(id, ShouldEqual, compare)
}

func (t *F) createIdQuery(id string) map[string]interface{} {
	return map[string]interface{}{"_id": bson.ObjectIdHex(id)}
}

func (t *F) createStringIdQuery(id string) map[string]interface{} {
	return map[string]interface{}{"_id": id}
}
