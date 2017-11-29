package api

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
)

// DownloadSource represents one file to upload.
//
// It is only valid to set one of (Writer, Path).
// If Path is set, it will be written to disk using os.Create.
type DownloadSource struct {
	Writer io.WriteCloser
	Path   string
}

func CreateDownloadSourceFromFilename(filename string) *DownloadSource {
	return &DownloadSource{Path: filename}
}

func (c *Client) Download(url string, progress chan<- int64, destination *DownloadSource) chan error {

	// Synchronous closure
	doDownload := func() error {
		// Open the writer based on destination path, if no writer was given.
		if destination.Writer == nil {
			if destination.Path == "" {
				return errors.New("Neither destination path nor writer was set in download source")
			}
			fileWriter, err := os.Create(destination.Path)
			if err != nil {
				return err
			}
			destination.Writer = fileWriter
		}
		defer destination.Writer.Close()

		req, err := c.New().Get(url).Request()
		if err != nil {
			return err
		}

		resp, err := c.Doer.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			// Needs robust handling for body & raw nils
			raw, _ := ioutil.ReadAll(resp.Body)
			return errors.New(string(raw))
		}

		if resp.Body == nil {
			return errors.New("Response body was empty")
		}

		// Pass response body through a ProgressReader which will report to the progress chan
		progressReader := NewProgressReader(resp.Body, progress)
		defer progressReader.Close()

		// Copy response
		var written int64
		written, err = io.Copy(destination.Writer, progressReader)

		// Verify that the written length was what was expected
		if resp.ContentLength != written {
			return errors.New("Response body was truncated")
		}

		return err
	}

	// Report result back to caller
	resultChan := make(chan error, 1)

	go func() {
		err := doDownload()
		resultChan <- err
	}()

	return resultChan
}

func (c *Client) DownloadSimple(url string, destination *DownloadSource) (chan int64, chan error) {

	progress := make(chan int64, 10)

	return progress, c.Download(url, progress, destination)
}
