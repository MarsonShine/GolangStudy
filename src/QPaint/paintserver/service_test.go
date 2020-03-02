package paintserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type idRet struct {
	ID string `json:"id"`
}

func newServer() *httptest.Server {
	// doc := NewDocument()
	// service := NewService(doc)
	// return httptest.NewServer(service)
	return nil
}

func Post(ret interface{}, url, body string) (err error) {
	b := strings.NewReader(body)
	resp, err := http.Post(url, "application/json", b)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if ret != nil {
		err = json.NewDecoder(resp.Body).Decode(ret)
	}
	return
}

func TestNewDrawing(t *testing.T) {
	ts := newServer()
	defer ts.Close()

	var ret idRet
	err := Post(&ret, ts.URL+"/drawings", "")
	if err != nil {
		t.Fatal("Post /drawings failed:", err)
	}
	if ret.ID != "10001" {
		t.Log("new drawing id:", ret.ID)
	}
}
