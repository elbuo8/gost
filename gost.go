package gost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	GistAPIURL   = "https://api.github.com"
	AcceptHeader = "application/vnd.github.v3+json"
)

type Gost struct {
	token      string
	Client     *http.Client
	GistAPIURL string
}

type Gist struct {
	Description string              `json:"description"`
	Public      bool                `json:"public"`
	Files       map[string]GistFile `json:"files"`
	Filename    string              `json:"filename,omitempty"`
}

type GistFile struct {
	Filename string `json:"-"`
	Content  string `json:"content"`
}

func timeoutHandler(network, address string) (net.Conn, error) {
	return net.DialTimeout(network, address, time.Duration(5*time.Second))
}

func NewGost(token string) *Gost {
	transport := http.Transport{
		Dial: timeoutHandler,
	}
	return &Gost{
		token: token,
		Client: &http.Client{
			Transport: &transport,
		},
		GistAPIURL: GistAPIURL,
	}
}

func (g *Gost) GetUser(user string) ([]byte, error) {
	return mitigateJSON(g.makeRequest("GET", "/users/"+user+"/gists", nil))
}

func (g *Gost) GetPublic() ([]byte, error) {
	return mitigateJSON(g.makeRequest("GET", "/gists/public", nil))
}

func (g *Gost) GetStarred() ([]byte, error) {
	return mitigateJSON(g.makeRequest("GET", "/gists/starred", nil))
}

func (g *Gost) Get(id string) ([]byte, error) {
	return mitigateJSON(g.makeRequest("GET", "/gists/"+id, nil))
}

func (g *Gost) Create(description string, public bool, files ...*GistFile) ([]byte, error) {
	gist := Gist{
		Description: description,
		Public:      public,
		Files:       make(map[string]GistFile),
	}
	for i := 0; i < len(files); i++ {
		if files[i].Filename == "" {
			return nil, fmt.Errorf("Filename undefined for %+v GistFile", files[i])
		}
		gist.Files[files[i].Filename] = *files[i]
	}
	payload, _ := json.Marshal(gist)
	return mitigateJSON(g.makeRequest("POST", "/gists", bytes.NewReader(payload)))
}

func (g *Gost) Edit(id string, gist *Gist) ([]byte, error) {
	payload, _ := json.Marshal(gist)
	return mitigateJSON(g.makeRequest("PATCH", "/gists/"+id, bytes.NewReader(payload)))
}

func (g *Gost) ListCommits(id string) ([]byte, error) {
	return mitigateJSON(g.makeRequest("GET", "/gists/"+id+"/commits", nil))
}

func (g *Gost) Fork(id string) ([]byte, error) {
	return mitigateJSON(g.makeRequest("POST", "/gists/"+id+"/forks", nil))
}

func (g *Gost) ListForks(id string) ([]byte, error) {
	return mitigateJSON(g.makeRequest("GET", "/gists/"+id+"/forks", nil))
}

func (g *Gost) Star(id string) (bool, error) {
	resp, err := g.makeRequest("PUT", "/gists/"+id+"/star", nil)
	return mitigateCode(http.StatusNoContent, resp, err)
}

func (g *Gost) UnStar(id string) (bool, error) {
	resp, err := g.makeRequest("DELETE", "/gists/"+id+"/star", nil)
	return mitigateCode(http.StatusNoContent, resp, err)
}

func (g *Gost) CheckStar(id string) (bool, error) {
	resp, err := g.makeRequest("GET", "/gists/"+id+"/star", nil)
	return mitigateCode(http.StatusNoContent, resp, err)
}

func (g *Gost) Delete(id string) (bool, error) {
	resp, err := g.makeRequest("DELETE", "/gists/"+id, nil)
	return mitigateCode(http.StatusNoContent, resp, err)
}

func (g *Gost) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(method, g.GistAPIURL+endpoint, body)
	header := http.Header{}
	header.Add("Accept", AcceptHeader)
	header.Add("Authorization", "token "+g.token)
	req.Header = header
	resp, err := g.Client.Do(req)
	switch {
	case err != nil:
		return nil, err
	case resp.StatusCode == http.StatusUnauthorized:
		return nil, fmt.Errorf("Invalid Token")
	}

	return resp, err
}

func mitigateCode(code int, resp *http.Response, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == code {
		return true, nil
	}
	return false, nil
}

func mitigateJSON(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
