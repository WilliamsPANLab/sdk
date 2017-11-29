package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"flywheel.io/sdk/api"
)

func init() {
	// Deterministically generating random numbers in parallel?
	// Sounds like a problem for another day.
	// Would probably use stack pointers, or ticket numbers, or something.
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var hexRunes = []rune("0123456789abcdef")

func RandStringOfLength(n int, runes []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func RandString() string {
	return RandStringOfLength(10, letterRunes)
}

func RandStringLower() string {
	return strings.ToLower(RandStringOfLength(10, letterRunes))
}

func RandHex() string {
	return RandStringOfLength(24, hexRunes)
}

// BEGIN: several variables lifted from smartystreets/assertions, because not exported :(
const (
	success         = ""
	shouldUseTimes  = "You must provide time instances as arguments to this assertion."
	needExactValues = "This assertion requires exactly %d comparison values (you provided %d)."
)

func need(needed int, expected []interface{}) string {
	if len(expected) != needed {
		return fmt.Sprintf(needExactValues, needed, len(expected))
	}
	return success
}

// END

const (
	shouldBeTimeEqual = "Expected: '%s'\nActual:   '%s'\n(Should be the same time, but they differed by %s)"
)

// Workaround for ShouldEqual and ShouldResemble being poor time.Time comparators.
// https://github.com/smartystreets/assertions/issues/15
func ShouldBeSameTimeAs(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}
	actualTime, firstOk := actual.(time.Time)
	expectedTime, secondOk := expected[0].(time.Time)

	if !firstOk || !secondOk {
		return shouldUseTimes
	}

	if !actualTime.Equal(expectedTime) {
		return fmt.Sprintf(shouldBeTimeEqual, actualTime, expectedTime, actualTime.Sub(expectedTime))
	}

	return success
}

// Helper function to compare a string against a regular expression
const (
	invalidRegex     = "Invalid match expression '%s': %s"
	regexBadMatch    = "Expected '%s' to match pattern: '%s'"
	shouldUseStrings = "You must provide string instances as arguments to this assertion."
)

func ShouldMatchRegex(actual interface{}, pattern ...interface{}) string {
	if fail := need(1, pattern); fail != success {
		return fail
	}

	actualString, firstOk := actual.(string)
	patternString, secondOk := pattern[0].(string)

	if !firstOk || !secondOk {
		return shouldUseStrings
	}

	matched, err := regexp.MatchString(patternString, actualString)
	if err != nil {
		return fmt.Sprintf(invalidRegex, patternString, err.Error())
	}
	if !matched {
		return fmt.Sprintf(regexBadMatch, actualString, patternString)
	}
	return success
}

func UploadSourceFromString(name, src string) *api.UploadSource {
	return &api.UploadSource{
		Reader: ioutil.NopCloser(bytes.NewBufferString(src)),
		Name:   name,
	}
}

type doFunc func(*http.Request) (*http.Response, error)
type httpResponseWrapperFunc func(*http.Response, error) (*http.Response, error)

func (f doFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func HttpResponseWrapper(client *api.Client, fn httpResponseWrapperFunc) *api.Client {
	doer := doFunc(func(req *http.Request) (*http.Response, error) {
		resp, err := client.Doer.Do(req)
		return fn(resp, err)
	})

	return &api.Client{
		doer,
		client.Sling.New().Doer(doer),
	}
}

func HttpResponseLengthSetter(len int64) httpResponseWrapperFunc {
	return func(resp *http.Response, err error) (*http.Response, error) {
		if resp != nil {
			resp.ContentLength = len
		}
		return resp, err
	}
}

// Buffer does not implement close; ioutil does not implement NopWriteCloser
type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error { return nil }
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}

func DownloadSourceToBuffer() (*bytes.Buffer, *api.DownloadSource) {
	buffer := new(bytes.Buffer)

	return buffer, &api.DownloadSource{
		Writer: NopWriteCloser(buffer),
	}
}

// TEMP

func Format(x interface{}) string {
	y, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(y)
}

func PrintFormat(x interface{}) {
	y, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(y))
	}
}
