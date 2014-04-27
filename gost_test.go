package gost

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	TOKEN = "MY TOKEN!"
)

var g = NewGost(TOKEN)

var fakeServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.Header.Get("Authorization") != "token "+TOKEN:
		http.Error(w, "Missing token", http.StatusUnauthorized)
	case req.Header.Get("Accept") != "application/vnd.github.v3+json":
		http.Error(w, "Missing Accept Header", http.StatusNotAcceptable)
	case req.Method == "GET" && req.Body != nil:
		http.Error(w, "GET should not have body", http.StatusBadRequest)
	case req.Method == "POST" || req.Method == "PATCH" && req.Body == nil:
		http.Error(w, "POST/PATCH should have body", http.StatusBadRequest)
	default:
		fmt.Fprintf(w, "Request Accepted!")
	}
}))

func TestNewGost(t *testing.T) {
	gost := NewGost(TOKEN)
	if gost == nil {
		t.Errorf("New Gost should not return nil")
	}
}

func TestGetUser(t *testing.T) {
	g.GistAPIURL = fakeServer.URL
	_, err := g.GetUser("defunkt")
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestGetPublic(t *testing.T) {
	g.GistAPIURL = fakeServer.URL
	_, err := g.GetPublic()
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestGetStarred(t *testing.T) {
	g.GistAPIURL = fakeServer.URL
	_, err := g.GetStarred()
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestGet(t *testing.T) {
	g.GistAPIURL = fakeServer.URL
	_, err := g.Get("2")
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestCreate(t *testing.T) {
	g.GistAPIURL = fakeServer.URL
	file1 := &GistFile{Filename: "Test.go", Content: "package go"}
	file2 := &GistFile{Filename: "Test2.go", Content: "package main"}
	_, err := g.Create("TESTTTT", false, file1, file2)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestCreateMissingName(t *testing.T) {
	g.GistAPIURL = fakeServer.URL
	file1 := &GistFile{Content: "package go"}
	file2 := &GistFile{Filename: "Test2.go", Content: "package main"}
	_, err := g.Create("TESTTTT", false, file1, file2)
	if err == nil {
		t.Errorf("Create must fail if filename is not provided")
	}
}

func TestEdit(t *testing.T) {
	g := NewGost(TOKEN)
	gist := &Gist{
		Description: "Awesome stuff",
		Filename:    "MY FILE!",
		Public:      false,
		Files:       make(map[string]GistFile),
	}
	gist.Files["myfile.go"] = GistFile{
		Filename: "Test2.go",
		Content:  "package main"}
	_, err := g.Edit("2", gist)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestListCommits(t *testing.T) {
	g.GistAPIURL = fakeServer.URL
	_, err := g.ListCommits("defunkt")
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestFork(t *testing.T) {
	g.GistAPIURL = fakeServer.URL
	_, err := g.Fork("2")
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestListForks(t *testing.T) {
	g.GistAPIURL = fakeServer.URL
	_, err := g.ListForks("defunkt")
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestFailedRequest(t *testing.T) {
	gost := NewGost("Fail Token")
	gost.GistAPIURL = fakeServer.URL
	_, err := gost.GetPublic()
	if err == nil {
		t.Errorf("Request should fail on incorrect token")
	}
}

func TestStopServers(t *testing.T) {
	fakeServer.Close()
}
