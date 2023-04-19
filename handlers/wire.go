package handlers

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/handlers/share/gallery"
	"github.com/steebchen/keskin-api/handlers/share/review"
	"github.com/steebchen/keskin-api/handlers/template/index"
	"github.com/steebchen/keskin-api/handlers/template/webmanifest"
)

var ProviderSet = wire.NewSet(
	review.New,
	gallery.New,
	index.New,
	webmanifest.New,
)
