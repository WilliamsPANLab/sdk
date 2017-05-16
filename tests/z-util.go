package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
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

func UploadSourceFromString(name, src string) *api.UploadSource {
	return &api.UploadSource{
		Reader: ioutil.NopCloser(bytes.NewBufferString(src)),
		Name:   name,
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
