package tests

import (
	. "github.com/smartystreets/assertions"

	"flywheel.io/sdk/api"
	"net/http"
)

func (t *F) TestCreateDownloadSourceFromFilenames() {
	source := api.CreateDownloadSourceFromFilename("one.txt")

	t.So(source.Path, ShouldEqual, "one.txt")
}

func (t *F) TestBadDownloads() {
	// Invalid download source
	source := &api.DownloadSource{}
	_, result := t.DownloadSimple("", source)
	t.So((<-result).Error(), ShouldEqual, "Neither destination path nor writer was set in download source")

	// Nonexistant download path
	source = &api.DownloadSource{Path: "/dev/null/does-not-exist"}
	_, result = t.DownloadSimple("", source)
	t.So((<-result).Error(), ShouldStartWith, "open /dev/null/does-not-exist: ")

	// Bad download url
	buffer, source := DownloadSourceToBuffer()
	_, result = t.DownloadSimple("not-an-endpoint", source)

	// Could improve this in the future
	err := <-result
	t.So(err.Error(), ShouldMatchRegex, "\\{\"status_code\": 404, \"message\": \"The resource could not be found.\"\\, \"request_id\": \"[^\"]+\"}")
	t.So(buffer.String(), ShouldEqual, "")
}

func (t *F) TestBadUrlDownload() {
	// Try with invalid project
	url, _, err := t.GetProjectDownloadUrl("not-a-project", "not-a-file")

	t.So(url, ShouldEqual, "")
	t.So(err, ShouldNotBeNil)
	t.So(err.Error(), ShouldEqual, "(404) The resource could not be found.")
}

func (t *F) TestTruncatedDownloads() {
	// Create test project, and upload text
	_, projectId := t.createTestProject()

	poem := "Surely some revelation is at hand;"
	t.uploadText(t.UploadToProject, projectId, "yeats.txt", poem)

	buffer, dest := DownloadSourceToBuffer()

	// Wrap the response handler
	client := HttpResponseWrapper(t.Client, HttpResponseLengthSetter(100))
	// Ignoring the progress channel result because a mismatch is expected
	_, result := client.DownloadFromProject(projectId, "yeats.txt", dest)

	err := <-result
	t.So(err, ShouldNotBeNil)
	t.So(err.Error(), ShouldEqual, "Response body was truncated")
	t.So(buffer.String(), ShouldEqual, poem)
}

// Given an download function, container ID, filename, and content - download & check content
func (t *F) downloadText(fn func(string, string, *api.DownloadSource) (chan int64, chan error), id, filename, text string) {
	buffer, dest := DownloadSourceToBuffer()
	progress, resultChan := fn(id, filename, dest)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, int64(len(text)))
	t.So(<-resultChan, ShouldBeNil)
	t.So(buffer.String(), ShouldEqual, text)
}

// Given a download function, container ID, filename, and content - generate a download ticket, download without authorization, and check content
func (t *F) downloadTextWithTicket(fn func(string, string) (string, *http.Response, error), id, filename, text string) {
	buffer, dest := DownloadSourceToBuffer()
	downloadUrl, _, err := fn(id, filename)

	t.So(err, ShouldBeNil)
	t.So(downloadUrl, ShouldNotEqual, "")

	noAuthClient := api.Client{
		Doer:  t.Doer,
		Sling: t.Sling.New().Set("Authorization", ""),
	}

	progress, resultChan := noAuthClient.DownloadSimple(downloadUrl, dest)

	// Last update should be the full string length.
	t.checkProgressChanEndsWith(progress, int64(len(text)))
	t.So(<-resultChan, ShouldBeNil)
	t.So(buffer.String(), ShouldEqual, text)
}
