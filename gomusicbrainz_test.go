package gomusicbrainz

// TODO use testdata from https://github.com/metabrainz/mmd-schema/tree/master/test-data/valid

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/michiwend/golang-pretty"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *WS2Client
)

// Init multiplexer and httptest server.
func setupHTTPTesting() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = NewWS2Client(
		server.URL,
		"Application Name",
		"Version",
		"Contact",
	)

	// NOTE this fixes testing since the test server does not listen on /ws/2
	client.WS2RootURL.Path = ""
}

// serveTestFile responses to the http client with content of a test file
// located in ./testdata.
func serveTestFile(endpoint string, testfile string, t *testing.T) {
	t.Log("Handling endpoint", endpoint)
	t.Log("Serving test file", testfile)

	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		t.Log("GET request was:", r.URL.String())

		http.ServeFile(w, r, path.Join("./testdata", testfile))
	})
}

// pretty prints a diff.
func requestDiff(want, returned interface{}) string {
	out := "\n"

	for _, diff := range pretty.Diff(want, returned) {
		out += fmt.Sprintln("difference in", diff)
	}
	return out
}
