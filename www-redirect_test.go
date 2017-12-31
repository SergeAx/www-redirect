package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type dnsNameTest struct {
	name   string
	result bool
}

var dnsNameTests = []dnsNameTest{
	// RFC 2181, section 11.
	{"_xmpp-server._tcp.google.com", true},
	{"foo.com", true},
	{"1foo.com", true},
	{"26.0.0.73.com", true},
	{"fo-o.com", true},
	{"fo1o.com", true},
	{"foo1.com", true},
	{"a.b..com", false},
	{"a.b-.com", false},
	{"a.b.com-", false},
	{"a.b..", false},
	{"b.com.", true},
}

func TestIsDomainName(t *testing.T) {
	for _, pair := range dnsNameTests {
		if isDomainName(pair.name) != pair.result {
			t.Errorf("isDomainName(%q) = %v; want %v", pair.name, !pair.result, pair.result)
		}
	}
}

type transformDomainTest struct {
	name   string
	result string
	status int
}

var transformDomainTests = []transformDomainTest{
	{"example.com", "www.example.com", 0},
	{"EXAMPLE.com", "www.example.com", 0},
	{"www.example.com", "", http.StatusNotFound},
	{"www.-example.com", "", http.StatusBadRequest},
	{"a", "www.a", 0},
}

func TestTransformDomain(t *testing.T) {
	for _, test := range transformDomainTests {
		result, status := transformDomain(test.name)
		if result == test.result && status == test.status {
			continue
		}
		t.Errorf("transformDomain(%s) = (%s, %d); want (%s, %d)", test.name, result, status, test.result, test.status)
	}
}

func TestRedirectHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/test-uri", nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Host", "e.com")

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectHandler)
	handler.ServeHTTP(recorder, request)

	status := recorder.Code
	wantedStatus := http.StatusMovedPermanently
	if status != wantedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, wantedStatus)
	}

	location := recorder.HeaderMap.Get("Location")
	wantedLocation := "http://www.e.com/test-uri"
	if location != wantedLocation {
		t.Errorf("handler returned wrong location: got %v want %v",
			location, wantedLocation)
	}

	server := recorder.HeaderMap.Get("Server")
	wantedServer := fmt.Sprintf("%s (%s)", serverSoftware, serverVersion)
	if server != wantedServer {
		t.Errorf("handler returned wrong server header: got %v want %v",
			server, wantedServer)
	}

	//TODO: test 404 and 400 statuses
}
