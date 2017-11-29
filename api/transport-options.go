package api

import (
	"io"

	"github.com/dghubble/sling"
)

// SdkDebugKey is the environment variable used to control request debug logging.
// If enabled, the SDK will log an HTTP/1.1 representation of all requests made to stdout.
const SdkDebugKey = "SdkDebug"

// Client is an http and sling client capable of making flywheel requests.
type Client struct {
	sling.Doer
	*sling.Sling
}

type ApiKeyClientOption func(*ApiKeyClientOptions)
type ApiKeyClientOptions struct {
	// Skip SSL verification
	InsecureSkipVerify bool

	// Use plaintext (HTTP) transport
	InsecureUsePlaintext bool

	// A writer to send debug request bodies to, if any
	DebugWriter io.Writer
}

var DefaultApiKeyClientOptions = ApiKeyClientOptions{
	InsecureSkipVerify:   false,
	InsecureUsePlaintext: false,
}

// Specify that the ApiKeyClient should not verify SSL connections.
// Should only be used for development.
var InsecureNoSSLVerification ApiKeyClientOption

// Specify that the ApiKeyClient should use a plaintext HTTP transport.
// Should only be used for development.
var InsecureUsePlaintext ApiKeyClientOption

// Specify that the ApiKeyClient should log all request bodies to the specified Writer.
// See DebugTransport for details.
func DebugLogRequests(w io.Writer) ApiKeyClientOption {
	return func(o *ApiKeyClientOptions) {
		o.DebugWriter = w
	}
}

func init() {
	InsecureNoSSLVerification = func(o *ApiKeyClientOptions) {
		o.InsecureSkipVerify = true
	}

	InsecureUsePlaintext = func(o *ApiKeyClientOptions) {
		o.InsecureUsePlaintext = true
	}
}
