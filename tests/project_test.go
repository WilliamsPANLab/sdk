package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestProjects() {
	groupId := t.createTestGroup()

	projectName := RandString()
	project := &api.Project{
		Name:        projectName,
		GroupId:     groupId,
		Description: "This is a description",
		Info: map[string]interface{}{
			"some-key": 37,
		},
	}

	// Add
	projectId, _, err := t.AddProject(project)
	t.So(err, ShouldBeNil)

	// Get
	rProject, _, err := t.GetProject(projectId)
	t.So(err, ShouldBeNil)
	t.So(rProject.Id, ShouldEqual, projectId)
	t.So(rProject.Name, ShouldEqual, project.Name)
	t.So(rProject.Description, ShouldEqual, project.Description)
	t.So(rProject.Info, ShouldContainKey, "some-key")
	t.So(rProject.Info["some-key"], ShouldEqual, 37)
	now := time.Now()
	t.So(*rProject.Created, ShouldHappenBefore, now)
	t.So(*rProject.Modified, ShouldHappenBefore, now)

	// Get all
	projects, _, err := t.GetAllProjects()
	t.So(err, ShouldBeNil)
	// workaround: all-container endpoints skip some fields, single-container does not. this sets up the equality check
	rProject.Files = nil
	rProject.Notes = nil
	rProject.Tags = nil
	rProject.Info = nil
	t.So(projects, ShouldContain, rProject)

	// Modify
	newName := RandString()
	projectMod := &api.Project{
		Name: newName,
		Info: map[string]interface{}{
			"another-key": 52,
		},
	}
	_, err = t.ModifyProject(projectId, projectMod)
	t.So(err, ShouldBeNil)
	changedProject, _, err := t.GetProject(projectId)
	t.So(changedProject.Name, ShouldEqual, newName)
	t.So(changedProject.Info, ShouldContainKey, "some-key")
	t.So(changedProject.Info, ShouldContainKey, "another-key")
	t.So(changedProject.Info["another-key"], ShouldEqual, 52)
	t.So(*changedProject.Created, ShouldBeSameTimeAs, *rProject.Created)
	t.So(*changedProject.Modified, ShouldHappenAfter, *rProject.Modified)

	// Notes, tags
	message := "This is a note"
	_, err = t.AddProjectNote(projectId, message)
	t.So(err, ShouldBeNil)
	tag := "example-tag"
	_, err = t.AddProjectTag(projectId, tag)
	t.So(err, ShouldBeNil)

	// Check
	rProject, _, err = t.GetProject(projectId)
	t.So(err, ShouldBeNil)
	t.So(rProject.Notes, ShouldHaveLength, 1)
	t.So(rProject.Notes[0].Text, ShouldEqual, message)
	t.So(rProject.Tags, ShouldHaveLength, 1)
	t.So(rProject.Tags[0], ShouldEqual, tag)

	// Delete
	_, err = t.DeleteProject(projectId)
	t.So(err, ShouldBeNil)
	projects, _, err = t.GetAllProjects()
	t.So(err, ShouldBeNil)
	t.So(projects, ShouldNotContain, rProject)
}

func (t *F) TestProjectFiles() {
	groupId := t.createTestGroup()

	project := &api.Project{Name: RandString(), GroupId: groupId}
	projectId, _, err := t.AddProject(project)
	t.So(err, ShouldBeNil)

	poem := "The ceremony of innocence is drowned;"
	t.uploadText(t.UploadToProject, projectId, "yeats.txt", poem)

	rProject, _, err := t.GetProject(projectId)
	t.So(err, ShouldBeNil)
	t.So(rProject.Files, ShouldHaveLength, 1)
	t.So(rProject.Files[0].Name, ShouldEqual, "yeats.txt")
	t.So(rProject.Files[0].Size, ShouldEqual, 37)
	t.So(rProject.Files[0].Mimetype, ShouldEqual, "text/plain")

	// Download the same file and check content
	t.downloadText(t.DownloadFromProject, projectId, "yeats.txt", poem)

	// Bundling: test file attributes
	t.So(rProject.Files[0].Modality, ShouldEqual, "")
	t.So(rProject.Files[0].Measurements, ShouldHaveLength, 0)
	t.So(rProject.Files[0].Type, ShouldEqual, "text")

	_, response, err := t.ModifyProjectFile(projectId, "yeats.txt", &api.FileFields{
		Modality:     "MR",
		Measurements: []string{"functional"},
		Type:         "dicom",
	})
	t.So(err, ShouldBeNil)

	// Check that no jobs were triggered and attrs were modified
	t.So(response.JobsTriggered, ShouldEqual, 0)

	rProject, _, err = t.GetProject(projectId)
	t.So(err, ShouldBeNil)
	t.So(rProject.Files[0].Modality, ShouldEqual, "MR")
	t.So(rProject.Files[0].Measurements, ShouldHaveLength, 1)
	t.So(rProject.Files[0].Measurements[0], ShouldEqual, "functional")
	t.So(rProject.Files[0].Type, ShouldEqual, "dicom")

	// Test file info
	t.So(rProject.Files[0].Info, ShouldBeEmpty)
	_, err = t.ReplaceProjectFileInfo(projectId, "yeats.txt", map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	})
	t.So(err, ShouldBeNil)
	_, err = t.SetProjectFileInfo(projectId, "yeats.txt", map[string]interface{}{
		"c": 5,
	})

	rProject, _, err = t.GetProject(projectId)
	t.So(err, ShouldBeNil)
	t.So(rProject.Files[0].Info["a"], ShouldEqual, 1)
	t.So(rProject.Files[0].Info["b"], ShouldEqual, 2)
	t.So(rProject.Files[0].Info["c"], ShouldEqual, 5)
	t.So(rProject.Files[0].Info["d"], ShouldEqual, 4)

	_, err = t.DeleteProjectFileInfoFields(projectId, "yeats.txt", []string{"c", "d"})
	t.So(err, ShouldBeNil)

	rProject, _, err = t.GetProject(projectId)
	t.So(err, ShouldBeNil)
	t.So(rProject.Files[0].Info["a"], ShouldEqual, 1)
	t.So(rProject.Files[0].Info["b"], ShouldEqual, 2)
	t.So(rProject.Files[0].Info["c"], ShouldBeNil)
	t.So(rProject.Files[0].Info["d"], ShouldBeNil)

	_, err = t.ReplaceProjectFileInfo(projectId, "yeats.txt", map[string]interface{}{})
	rProject, _, err = t.GetProject(projectId)
	t.So(err, ShouldBeNil)
	t.So(rProject.Files[0].Info, ShouldBeEmpty)
}

func (t *F) createTestProject() (string, string) {
	groupId := t.createTestGroup()

	project := &api.Project{
		Name:    RandString(),
		GroupId: groupId,
	}
	projectId, _, err := t.AddProject(project)
	t.So(err, ShouldBeNil)

	return groupId, projectId
}
