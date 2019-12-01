package cors

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	router := AddRoutes(mux.NewRouter())
	router.
		Methods("GET").
		Path("/").
		HandlerFunc(notfound)

	// create a basic route
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	// validate the route returns notfound from below
	if w.Code != http.StatusNotFound {
		t.Fatal("Expected Notfound")
	}

	// Validate that the OPTIONS request for that route returns the expected results
	r = httptest.NewRequest("OPTIONS", "/", nil)
	r.Header.Add("Access-Control-Request_anything", "mydata")
	r.Header.Add("Random-Header", "randomData")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	// OPTIONS returns OK
	if w.Code != http.StatusOK {
		t.Fatal("got response code", w.Code)
	}
	// Headers are correct
	if m := w.Header().Get(allowMethod); m != allowedMethod {
		t.Fatal("Allowed methods do not match", m, allowedMethod)
	}
}

func notfound(w http.ResponseWriter, r *http.Request) {
	// dummy return value for GET
	http.Error(w, "Expected Notfound", http.StatusNotFound)
}
