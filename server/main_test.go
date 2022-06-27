package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// Helper function to process a request and test its response
func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

func TestHelloWorld(t *testing.T) {
	r := gin.Default()

	r.GET("/", HelloWorld)

	req, _ := http.NewRequest("GET", "/", nil)

	testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK
		if !statusOK {
			t.Errorf("wrong status expected %v got %v", http.StatusOK, w.Code)
		}
		p, err := ioutil.ReadAll(w.Body)
		bodyOK := err == nil && strings.Index(string(p), "Hello World!") > 0
		if !bodyOK {
			t.Errorf("wrong payload expected %v got %v", "Hello World!", w.Body)
		}
		return statusOK && bodyOK
	})
}
