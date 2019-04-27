package web

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var updateGoldenFiles = flag.Bool("update", false, "update golden files")

// AssertHTTP performs an HTTP call to the passed handler and asserts the
// resulting status code matches the expectation.
//
// AssertHTTP returns the response header and body for further assertions.
func AssertHTTP(
	t *testing.T, handler http.HandlerFunc, method, url string, values url.Values, body io.Reader, code int,
) (http.Header, []byte) {
	req := httptest.NewRequest(method, url, body)
	req.URL.RawQuery = values.Encode()
	rr := httptest.NewRecorder()
	handler(rr, req)
	assert.Equalf(t, code, rr.Code, "expected status code %d; got %d", code, rr.Code)
	return rr.Header(), rr.Body.Bytes()
}

// MustParseURL parses rawurl or fails the test if parsing fails.
func MustParseURL(t *testing.T, rawurl string) *url.URL {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		t.Fatalf("parse url %s: %v", rawurl, err)
	}
	return parsed
}

// NewURLValues encodes the passed kvs as url query parameters. It fails
// the amount of kvs is not even.
func NewURLValues(t *testing.T, kvs ...string) url.Values {
	if len(kvs)%2 != 0 {
		t.Fatalf("NewURLValues: even number of kvs expected")
	}
	values := url.Values{}
	for i := 0; i < len(kvs); i += 2 {
		values.Add(kvs[i], kvs[i+1])
	}
	return values
}

// NewHTTPHeader encodes the passed kvs as http header. It fails
// the amount of kvs is not even.
func NewHTTPHeader(t *testing.T, kvs ...string) http.Header {
	if len(kvs)%2 != 0 {
		t.Fatalf("NewHTTPHeader: even number of kvs expected")
	}
	header := http.Header{}
	for i := 0; i < len(kvs); i += 2 {
		header.Add(kvs[i], kvs[i+1])
	}
	return header
}

// AssertHTTPGoldenFile compares the passed header and body against the golden
// file for the test t.
func AssertHTTPGoldenFile(t *testing.T, header http.Header, body []byte) bool {
	testdata := filepath.Join("testdata", t.Name())
	if err := os.MkdirAll(testdata, 0744); err != nil {
		t.Fatal(err)
	}
	headerPath := filepath.Join(testdata, "header.golden.json")
	bodyPath := filepath.Join(testdata, "body.golden.html")
	if *updateGoldenFiles {
		writeHeaderAsJSON(t, headerPath, header)
		if err := ioutil.WriteFile(bodyPath, body, 0644); err != nil {
			t.Fatal(err)
		}
	}
	goldenHeader := readHeaderFromJSON(t, headerPath)
	goldenBody, err := ioutil.ReadFile(bodyPath)
	if err != nil {
		t.Fatal(err)
	}
	return assert.Equal(t, goldenHeader, header) &&
		assert.Equal(t, string(goldenBody), string(body))
}

func writeHeaderAsJSON(t *testing.T, headerPath string, header http.Header) {
	bs, err := json.MarshalIndent(header, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(headerPath, bs, 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func readHeaderFromJSON(t *testing.T, headerPath string) http.Header {
	headerFile, err := os.Open(headerPath)
	if err != nil {
		t.Fatal(err)
	}
	defer headerFile.Close()
	header := http.Header{}
	err = json.NewDecoder(headerFile).Decode(&header)
	if err != nil {
		t.Fatal(err)
	}
	return header
}
