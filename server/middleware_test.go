package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// import (
// 	"log"
// 	"net/http"
// 	"strings"
// )

func MockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestRestrictListing(t *testing.T) {
	testCases := []struct {
		desc     string
		url      string
		expected int
	}{
		{
			desc:     "Should not stop urls without / ending",
			url:      "/not-dir",
			expected: http.StatusOK,
		},
		{
			desc:     "Should stop urls with / ending",
			url:      "/dir/",
			expected: http.StatusNotFound,
		},
	}
	testHandler := http.HandlerFunc(MockHandler)
	handler := RestrictListing(testHandler)

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)
			if status := res.Code; status != test.expected {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, test.expected)
			}
		})
	}

}

func TestGetMethodOnly(t *testing.T) {
	testCases := []struct {
		desc     string
		method   string
		expected int
	}{
		{
			desc:     "Should accept GET request",
			method:   http.MethodGet,
			expected: http.StatusOK,
		},
		{
			desc:     "Should refuse POST reqeuse",
			method:   http.MethodPost,
			expected: http.StatusMethodNotAllowed,
		},
		{
			desc:     "Should refuse DELETE request",
			method:   http.MethodDelete,
			expected: http.StatusMethodNotAllowed,
		},
		{
			desc:     "Should refuse PUT request",
			method:   http.MethodPut,
			expected: http.StatusMethodNotAllowed,
		},
		{
			desc:     "Should refuse HEAD request",
			method:   http.MethodHead,
			expected: http.StatusMethodNotAllowed,
		},
		{
			desc:     "Should refuse OPTIONS request",
			method:   http.MethodOptions,
			expected: http.StatusMethodNotAllowed,
		},
	}
	testHandler := http.HandlerFunc(MockHandler)
	handler := GetMethodOnly(testHandler)

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest(test.method, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)
			if status := res.Code; status != test.expected {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, test.expected)
			}
		})
	}
}
