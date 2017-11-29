package api

import (
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// UploadSource represents one file to upload.
//
// It is only valid to set one of (Reader, Path).
// If Path is set, it will be read off disk using os.Open.
//
// If Name is not set, then filepath.Base(Path) will be used.
type UploadSource struct {
	Name string

	Reader io.ReadCloser
	Path   string
}

// Bundle an http response and error together for returning over a channel
type UploadResponse struct {
	Response *http.Response
	Error    error
}

func CreateUploadSourceFromFilenames(filenames ...string) []*UploadSource {
	var sources = make([]*UploadSource, len(filenames))

	for x, filename := range filenames {
		sources[x] = &UploadSource{Path: filename}
	}

	return sources
}

// Write a set of UploadSources to a multipart writer, reporting progress to a ProgressReader.
func writeUploadSources(writer *multipart.Writer, reader *ProgressReader, metadata []byte, files []*UploadSource) error {
	defer func() {
		writer.Close()
		reader.Close()
	}()

	// Add metadata, if any
	if len(metadata) > 0 {
		mWriter, err := writer.CreateFormField("metadata")
		if err != nil {
			return err
		}
		_, err = mWriter.Write(metadata)
		if err != nil {
			return err
		}
	}

	for i, file := range files {
		// Name the file, if no name was given.
		if file.Name == "" {
			if file.Path == "" {
				return errors.New("Neither file name nor path was set in upload source")
			}
			file.Name = filepath.Base(file.Path)
		}

		// Open a file descriptor if this UploadSource was not already an open reader
		if file.Reader == nil {
			fileReader, err := os.Open(file.Path)
			file.Reader = fileReader

			if err != nil {
				return err
			}
		}
		defer file.Reader.Close()

		// Report progress of the uploads, not of the encoded stream.
		// Upload progress of metadata and preamble will not be reported.
		reader.SetReader(file.Reader)

		// Create a form name for this file.
		// If there's only one file, don't add an index.
		// It might be valid to upload without this check. Worth testing.
		formTitle := "file"
		if len(files) > 1 {
			formTitle = strings.Join([]string{"file", strconv.Itoa(i + 1)}, "")
		}

		// Create a form entry for this file
		fileWriter, err := writer.CreateFormFile(formTitle, file.Name)
		if err != nil {
			return err
		}

		// Copy the file
		_, err = io.Copy(fileWriter, reader)
		if err != nil && err != io.EOF {
			return err
		}

		file.Reader.Close()
	}

	return nil
}

// Fire an upload with a given url, reader, and content type.
func (c *Client) sendUploadRequest(url string, reader io.ReadCloser, contentType string) (*http.Response, error) {
	req, err := c.New().Post(url).
		Body(reader).
		Set("Content-Type", contentType).
		Request()

	if err != nil {
		return nil, err
	}

	resp, err := c.Doer.Do(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != 200 {
		// Needs robust handling for body & raw nils
		raw, _ := ioutil.ReadAll(resp.Body)

		return resp, errors.New(string(raw))
	}

	return resp, err
}

// Upload will send a set of UploadSources to url, reporting uploaded bytes to progress if set.
// Upload will not block sending to progress.
//
// Depending on the URL, metadata may be required, or only one file may be allowed at a time.
// It is generally a good idea to use a purpose-specific upload method.
func (c *Client) Upload(url string, metadata []byte, progress chan<- int64, files []*UploadSource) chan error {

	// Form data is written from one goroutine to another
	reader, writer := io.Pipe()
	multipartWriter := multipart.NewWriter(writer)
	contentType := multipartWriter.FormDataContentType()

	// Wrap the pipe in a ProgressReader.
	progressReader := NewProgressReader(nil, progress)

	// Shared memory for results, protected by a waitgroup. Simpler (but more dangerous) than channels.
	var writeError error
	var uploadError error
	var response *http.Response
	var wg sync.WaitGroup
	wg.Add(2)

	// Report result back to caller
	resultChan := make(chan error, 1)

	// Stream multipart encoding
	go func() {
		writeError = writeUploadSources(multipartWriter, progressReader, metadata, files)
		writer.Close()
		wg.Done()
	}()

	// Send encoded body to server, await completion, report
	go func() {
		response, uploadError = c.sendUploadRequest(url, reader, contentType)
		wg.Done()
	}()

	// Wait for both to complete, and report back
	go func() {
		wg.Wait()

		// Encoding & local-IO errors take precedence over network errors.
		// Could combine the two if both are set. Eh.
		if writeError != nil {
			resultChan <- writeError
		} else {
			resultChan <- uploadError
		}
	}()

	return resultChan
}

// UploadSimple is a convenience wrapper around Upload.
// It creates the progress channel and UploadSource array for you.
func (c *Client) UploadSimple(url string, metadata []byte, files ...*UploadSource) (chan int64, chan error) {

	progress := make(chan int64, 10)

	return progress, c.Upload(url, metadata, progress, files)
}
