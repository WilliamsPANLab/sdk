package api

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/dghubble/sling"
)

type RootModeQueryStruct struct {
	Root bool `url:"root,omitempty"`
}

// NewApiKeyClient creates a Client with the given API key and options.
// Passing a key with an invalid format will panic.
func NewApiKeyClient(apiKey string, options ...ApiKeyClientOption) *Client {

	// If the debug environment variable is set, add the debug transport.
	// This is added to the beginning of the options array, in case it is overridden later.
	_, debug := os.LookupEnv(SdkDebugKey)
	if debug {
		options = append([]ApiKeyClientOption{DebugLogRequests(os.Stderr)}, options...)
	}

	// Load all configuration options into a config struct
	config := DefaultApiKeyClientOptions
	for _, x := range options {
		x(&config)
	}

	host, port, key, err := ParseApiKey(apiKey)
	if err != nil {
		panic(err)
	}

	// Load TLS configuration into a transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.InsecureSkipVerify,
		},
	}

	var hc *http.Client

	// Create an HTTP client, adding the debug transport if specified
	if config.DebugWriter == nil {
		hc = &http.Client{
			Transport: tr,
		}
	} else {
		debugT := &DebugTransport{
			Transport: tr,
			Writer:    config.DebugWriter,
		}
		hc = debugT.Client()
	}

	protocol := "https"
	if config.InsecureUsePlaintext {
		protocol = "http"
	}

	// Create a sling client, which is used for most server interactions
	sc := sling.New().
		Base(protocol+"://"+host+":"+strconv.Itoa(port)+"/").
		Set("Authorization", "scitran-user "+key).
		Set("User-Agent", "Flywheel SDK").
		Path("api/").
		Client(hc)

	// Set root query param for all requests, if requested
	if config.EnableRoot {
		var params = &RootModeQueryStruct{Root: true}
		sc.QueryStruct(params)
	}

	return &Client{
		Doer:  hc,
		Sling: sc,
	}
}

// DebugTransport prints its raw requests and responses to a writer, including bodies.
type DebugTransport struct {

	// All requests made with this transport will be written to Writer.
	// It will default to os.Stderr if nil.
	Writer io.Writer

	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	Transport http.RoundTripper
}

var header = []byte("\n--------------------\n---- BEGIN HTTP ----\n--------------------\n")
var middle = []byte("\n--------------------\n-- BEGIN RESPONSE --\n--------------------\n")
var footer = []byte("\n--------------------\n----- END HTTP -----\n--------------------\n")

// RoundTrip implements the RoundTripper interface.
func (t *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	// Effort in this made in this function, to:
	//
	// 1) Write debug output only once, to be concurrency-friendly
	// 2) Write as much debug output as possible when errors occur
	// 3) Expose the most relevant error: return an HTTP error if there was both an HTTP and a debug error.
	//
	// For that reason, there's a bit of "if err is nil okay now we can use err" logic.

	msg, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		Println("ergh")
		return nil, err
	}

	msg = append(header, msg...)
	msg = append(msg, middle...)

	var resp *http.Response

	if t.Transport != nil {
		resp, err = t.Transport.RoundTrip(req)
	} else {
		resp, err = http.DefaultTransport.RoundTrip(req)
	}

	// Stop if request failed
	if err != nil {
		str := []byte("Request failed: " + err.Error())
		msg = append(msg, str...)

	} else {
		// If there is a valid response, dump that too
		if resp != nil {
			var msg2 []byte
			msg2, err = httputil.DumpResponse(resp, true)

			if err != nil {
				str := []byte("Response dump failed: " + err.Error())
				msg = append(msg, str...)
			} else {
				msg = append(msg, msg2...)
			}
		}
	}

	msg = append(msg, footer...)

	var err2 error
	if t.Writer == nil {
		_, err2 = os.Stderr.Write(msg)
	} else {
		_, err2 = t.Writer.Write(msg)
	}

	// Don't allow debug write error to overwrite original error
	if err == nil {
		err = err2
	}

	return resp, err
}

// Client returns an *http.Client which uses the DebugTransport.
func (t *DebugTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}
