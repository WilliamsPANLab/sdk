package tests

import (
	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
)

func (t *F) TestCreateUploadSourceFromFilenames() {
	sources := api.CreateUploadSourceFromFilenames("one.txt", "two.txt")

	t.So(sources[0].Path, ShouldEqual, "one.txt")
	t.So(sources[1].Path, ShouldEqual, "two.txt")
}

func (t *F) TestBadUploads() {
	// Invalid upload source
	source := &api.UploadSource{}
	_, result := t.UploadSimple("", nil, source)
	t.So((<-result).Error(), ShouldEqual, "Neither file name nor path was set in upload source")

	// Nonexistant upload path
	source = &api.UploadSource{Path: "/dev/null/does-not-exist"}
	_, result = t.UploadSimple("", nil, source)
	t.So((<-result).Error(), ShouldStartWith, "open /dev/null/does-not-exist: ")

	// Bad upload url
	poem := "Are full of passionate intensity."
	source = UploadSourceFromString("yeats.txt", poem)
	_, result = t.UploadSimple("not-an-endpoint", nil, source)

	// Could improve this in the future
	err := <-result
	t.So(err.Error(), ShouldMatchRegex, "\\{\"status_code\": 404, \"message\": \"The resource could not be found.\"\\, \"request_id\": \"[^\"]+\"}")
}

// Given an upload function, container ID, filename, and content - upload & check length
func (t *F) uploadText(fn func(string, ...*api.UploadSource) (chan int64, chan error), id, filename, text string) {
	src := UploadSourceFromString(filename, text)
	progress, resultChan := fn(id, src)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, int64(len(text)))
	t.So(<-resultChan, ShouldBeNil)
}
