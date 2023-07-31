package server

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/steebchen/keskin-api/api"
	"github.com/steebchen/keskin-api/handlers/share/gallery"
	"github.com/steebchen/keskin-api/handlers/share/review"
	"github.com/steebchen/keskin-api/handlers/template/index"
	"github.com/steebchen/keskin-api/handlers/template/webmanifest"
	"github.com/steebchen/keskin-api/lib/file"
)

func NewServeMux(api *api.Handler, reviewHandler *review.Handler, galleryHandler *gallery.Handler, indexHandler *index.Handler, webmanifestHandler *webmanifest.Handler, config *Config) *http.ServeMux {
	mux := &http.ServeMux{}

	mux.Handle("/api/playground", playground.Handler("GraphQL Playground", "/api/graphql"))
	mux.Handle("/api/graphql", api)
	mux.Handle("/share/review/", reviewHandler)
	mux.Handle("/share/gallery/", galleryHandler)
	mux.Handle("/template/index", indexHandler)
	mux.Handle("/template/webmanifest", webmanifestHandler)

	mux.Handle(file.BasePath, http.StripPrefix(file.BasePath, http.FileServer(http.Dir(config.DataFolder))))

	return mux
}
