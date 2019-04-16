package tc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/go-kit/kit/log"

	"github.com/jho4us/tc/test"

	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestTransport(t *testing.T) {
	rep, err := testCreateTestRepo(t)
	if err != nil {
		t.Fatal(err)
	}
	tcs := NewService(rep)
	var logger log.Logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	httpLogger := log.With(logger, "component", "http")

	h := MakeHandler(tcs, httpLogger)

	testTruncateAll(tcs, t)

	id := testCreateHandler(h, t)
	testPutHandler(id, h, t)
	testGetHandler(tcs, h, t, id)
	testDeleteHandler(id, h, t)
}

func testCreateHandler(h http.Handler, t *testing.T) test.ID {
	body := &Test{
		Name:    "Aaaa",
		Content: "Bbbb",
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(body)

	req, err := http.NewRequest("POST", "/tc/v1/tests", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		t.Errorf("expected to return error on null post: got %v",
			rr.Code)
	}

	req, err = http.NewRequest("POST", "/tc/v1/tests", strings.NewReader("z=abra"))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		t.Errorf("expected to return error on invalid post: got %v",
			rr.Code)
	}

	req, err = http.NewRequest("POST", "/tc/v1/tests", buf)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler failed: got %v",
			rr.Code)
	}
	var ct createTestResponse
	resp := rr.Body.String()
	err = json.Unmarshal([]byte(resp), &ct)
	if err != nil {
		t.Errorf("expected to return json but %v returned",
			resp)

	}
	return ct.ID
}

func testPutHandler(id test.ID, h http.Handler, t *testing.T) {
	body := &Test{
		ID:      string(id),
		Name:    "Ccc",
		Content: "Xxx",
	}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("/tc/v1/tests/%s", id), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		t.Errorf("expected to return error on null post: got %v",
			rr.Code)
	}

	req, err = http.NewRequest("POST", fmt.Sprintf("/tc/v1/tests/%s", id), strings.NewReader("z=abra"))
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code == http.StatusOK {
		t.Errorf("expected to return error on invalid post: got %v",
			rr.Code)
	}

	req, err = http.NewRequest("POST", fmt.Sprintf("/tc/v1/tests/%s", id), buf)
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("handler failed: got %v",
			rr.Code)
	}

}

func testDeleteHandler(id test.ID, h http.Handler, t *testing.T) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/tc/v1/tests/%s", id), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler failed: got %v",
			rr.Code)
	}
}

func testGetHandler(s Service, h http.Handler, t *testing.T, id test.ID) {
	tt := []struct {
		routePath  string
		shouldPass bool
	}{
		{"tests", true},
		{fmt.Sprintf("tests/%s", id), true},
		{"tests/nonexistent", false},
		{"heap", false},
	}

	for _, tc := range tt {
		path := fmt.Sprintf("/tc/v1/%s", tc.routePath)
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		// In this case, our MetricsHandler returns a non-200 response
		// for a route variable it doesn't know about.
		if rr.Code == http.StatusOK && !tc.shouldPass {
			t.Errorf("handler should have failed on routePath %s: got %v want %v",
				tc.routePath, rr.Code, http.StatusOK)
		}
		if rr.Code != http.StatusOK && tc.shouldPass {
			t.Errorf("handler should have succeeded on routePath %s: got %v want %v",
				tc.routePath, rr.Code, http.StatusOK)
		}
	}
}
