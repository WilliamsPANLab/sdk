package tests

import (
	"time"

	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestGroups() {
	groupId := RandStringLower()
	groupName := RandString()

	group := &api.Group{
		Id:   groupId,
		Name: groupName,
	}

	// Add
	rId, _, err := t.AddGroup(group)
	t.So(err, ShouldBeNil)
	t.So(rId, ShouldEqual, groupId)

	// Get
	savedGroup, _, err := t.GetGroup(groupId)
	t.So(err, ShouldBeNil)
	t.So(savedGroup.Id, ShouldEqual, group.Id)
	t.So(savedGroup.Name, ShouldEqual, group.Name)
	now := time.Now()
	t.So(*savedGroup.Created, ShouldHappenBefore, now)
	t.So(*savedGroup.Modified, ShouldHappenBefore, now)

	// Get all
	groups, _, err := t.GetAllGroups()
	t.So(err, ShouldBeNil)
	// workaround: all-container endpoints skip some fields, single-container does not. this sets up the equality check
	// savedGroup.Permissions = nil
	t.So(groups, ShouldContain, savedGroup)

	// Modify
	newName := RandString()
	groupMod := &api.Group{
		Name: newName,
	}
	_, err = t.ModifyGroup(groupId, groupMod)
	t.So(err, ShouldBeNil)
	changedGroup, _, err := t.GetGroup(groupId)
	t.So(changedGroup.Name, ShouldEqual, newName)
	t.So(*changedGroup.Created, ShouldBeSameTimeAs, *savedGroup.Created)

	// Tags
	tag := "example-tag-group"
	_, err = t.AddGroupTag(groupId, tag)
	t.So(err, ShouldBeNil)

	// Check
	rGroup, _, err := t.GetGroup(groupId)
	t.So(err, ShouldBeNil)
	t.So(rGroup.Tags, ShouldHaveLength, 1)
	t.So(rGroup.Tags[0], ShouldEqual, tag)
	t.So(*changedGroup.Modified, ShouldHappenAfter, *savedGroup.Modified)

	// Delete
	_, err = t.DeleteGroup(groupId)
	t.So(err, ShouldBeNil)
	groups, _, err = t.GetAllGroups()
	t.So(err, ShouldBeNil)
	t.So(groups, ShouldNotContain, savedGroup)
}

func (t *F) createTestGroup() string {
	groupId := RandStringLower() // conform to group ID regex
	groupName := RandString()

	_, _, err := t.AddGroup(&api.Group{Id: groupId, Name: groupName})
	t.So(err, ShouldBeNil)

	return groupId
}
