package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/steebchen/keskin-api/api"
	"github.com/steebchen/keskin-api/handlers/share/gallery"
	"github.com/steebchen/keskin-api/handlers/share/review"
	"github.com/steebchen/keskin-api/handlers/template/index"
	"github.com/steebchen/keskin-api/handlers/template/webmanifest"
)

type mockHandler struct{}

func (h *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "success")
	if err != nil {
		panic(err)
	}
}

func TestNewServeMux(t *testing.T) {
	mux := NewServeMux(&api.Handler{
		Next: &mockHandler{},
	}, &review.Handler{}, &gallery.Handler{}, &index.Handler{}, &webmanifest.Handler{}, &Config{})

	server := httptest.NewServer(mux)
	defer server.Close()

	graphql := request(t, server, "/api/graphql")
	assert.Equal(t, "success", graphql)

	playground := request(t, server, "/api/playground")
	assert.Equal(t, true, strings.HasPrefix(playground, "<!DOCTYPE html>"))
}

func request(t *testing.T, server *httptest.Server, path string) string {
	resp, err := http.Get(server.URL + path)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	result, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	return string(result)
}
