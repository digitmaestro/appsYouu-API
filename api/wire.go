package api

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/api/resolver/appointment"
	"github.com/steebchen/keskin-api/api/resolver/branch"
	"github.com/steebchen/keskin-api/api/resolver/company"
	"github.com/steebchen/keskin-api/api/resolver/email_template"
	"github.com/steebchen/keskin-api/api/resolver/favorite"
	"github.com/steebchen/keskin-api/api/resolver/mutation"
	"github.com/steebchen/keskin-api/api/resolver/news"
	"github.com/steebchen/keskin-api/api/resolver/product"
	"github.com/steebchen/keskin-api/api/resolver/product_service_attribute"
	"github.com/steebchen/keskin-api/api/resolver/query"
	"github.com/steebchen/keskin-api/api/resolver/review"
	"github.com/steebchen/keskin-api/api/resolver/root"
	"github.com/steebchen/keskin-api/api/resolver/service"
	"github.com/steebchen/keskin-api/api/resolver/user"
)

var ProviderSet = wire.NewSet(
	query.ProviderSet,
	mutation.ProviderSet,

	appointment.ProviderSet,
	branch.ProviderSet,
	company.ProviderSet,
	user.ProviderSet,
	product.ProviderSet,
	service.ProviderSet,
	review.ProviderSet,
	favorite.ProviderSet,
	news.ProviderSet,
	email_template.ProviderSet,
	product_service_attribute.ProviderSet,

	root.ProviderSet,

	New,
)
