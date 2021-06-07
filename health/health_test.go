package health

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func testParseResponse(resp *http.Response) *CheckResponse {
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &CheckResponse{}
	}

	var r CheckResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return &CheckResponse{}
	}

	return &r
}

func TestHandler_OK(t *testing.T) {
	expected := &CheckResponse{
		ReleaseID: "42",
		ServiceID: "42",
		Version:   "42",
		Checks: map[string]*CheckResult{
			"name": {
				Output:        "privet",
				ObservedValue: float64(42),
				Status:        CheckStatusWarn,
			},
		},
		Status: CheckStatusPass,
	}

	ch := NewChecker(CheckerOptions{
		ReleaseID: "42",
		ServiceID: "42",
		Version:   "42",
	})
	r := &CheckResult{
		Output:        "privet",
		Status:        CheckStatusWarn,
		ObservedValue: 42,
	}
	ch.AddCallback("name", CheckCallback(func(_ context.Context) *CheckResult {
		return r
	}))
	srv := httptest.NewServer(http.HandlerFunc(NewHandler(ch, "someName")))
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL)
	if err != nil {
		t.Errorf("failed to perform healthcheck: %s", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	got := testParseResponse(resp)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got response %+v, expected %+v", got, expected)
	}

	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "application/json; charset=utf-8"
	if contentType != expectedContentType {
		t.Errorf("got Content-Type header: %s, expected: %s", contentType, expectedContentType)
	}
}

func TestHandler_Failed(t *testing.T) {
	expected := &CheckResponse{
		ReleaseID: "42",
		ServiceID: "42",
		Version:   "42",
		Output:    "privet",
		Checks: map[string]*CheckResult{
			"name": {
				Output:        "privet",
				ObservedValue: float64(42),
				Status:        CheckStatusFail,
			},
		},
		Status: CheckStatusFail,
	}

	ch := NewChecker(CheckerOptions{
		ReleaseID: "42",
		ServiceID: "42",
		Version:   "42",
	})
	r := &CheckResult{
		Output:        "privet",
		Status:        CheckStatusFail,
		ObservedValue: float64(42),
	}
	ch.AddCallback("name", func(_ context.Context) *CheckResult {
		return r
	})
	srv := httptest.NewServer(http.HandlerFunc(NewHandler(ch, "someOtherName")))
	defer srv.Close()

	resp, err := srv.Client().Get(srv.URL)
	if err != nil {
		t.Errorf("failed to perform healthcheck: %s", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	got := testParseResponse(resp)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got response %+v, expected %+v", got, expected)
	}
}
