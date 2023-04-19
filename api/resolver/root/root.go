package root

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
	"github.com/steebchen/keskin-api/api/resolver/service"
	"github.com/steebchen/keskin-api/api/resolver/user"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

type Resolver struct {
	Prisma                          *prisma.Client
	MutationResolver                *mutation.Mutation
	QueryResolver                   *query.Query
	UserResolver                    *user.User
	AppointmentResolver             *appointment.Appointment
	BranchResolver                  *branch.Branch
	CompanyResolver                 *company.Company
	ProductResolver                 *product.Product
	ServiceResolver                 *service.Service
	ReviewResolver                  *review.Review
	FavoriteResolver                *favorite.Favorite
	NewsResolver                    *news.News
	EmailTemplateResolver           *email_template.EmailTemplate
	ProductServiceAttributeResolver *product_service_attribute.ProductServiceAttribute
}

func New(
	client *prisma.Client,
	mutation *mutation.Mutation,
	query *query.Query,
	user *user.User,
	appointment *appointment.Appointment,
	branch *branch.Branch,
	company *company.Company,
	product *product.Product,
	service *service.Service,
	review *review.Review,
	favorite *favorite.Favorite,
	news *news.News,
	emailTemplate *email_template.EmailTemplate,
	productServiceAttribute *product_service_attribute.ProductServiceAttribute,
) *Resolver {
	return &Resolver{
		client,
		mutation,
		query,
		user,
		appointment,
		branch,
		company,
		product,
		service,
		review,
		favorite,
		news,
		emailTemplate,
		productServiceAttribute,
	}
}

func (r *Resolver) Mutation() gqlgen.MutationResolver {
	return r.MutationResolver
}

func (r *Resolver) Query() gqlgen.QueryResolver {
	return r.QueryResolver
}

func (r *Resolver) Appointment() gqlgen.AppointmentResolver {
	return r.AppointmentResolver
}

func (r *Resolver) Branch() gqlgen.BranchResolver {
	return r.BranchResolver
}

func (r *Resolver) Company() gqlgen.CompanyResolver {
	return r.CompanyResolver
}

func (r *Resolver) Customer() gqlgen.CustomerResolver {
	return r.UserResolver.CustomerResolver
}

func (r *Resolver) Employee() gqlgen.EmployeeResolver {
	return r.UserResolver.EmployeeResolver
}

func (r *Resolver) Manager() gqlgen.ManagerResolver {
	return r.UserResolver.ManagerResolver
}

func (r *Resolver) Administrator() gqlgen.AdministratorResolver {
	return r.UserResolver.AdministratorResolver
}

func (r *Resolver) Product() gqlgen.ProductResolver {
	return r.ProductResolver
}

func (r *Resolver) Service() gqlgen.ServiceResolver {
	return r.ServiceResolver
}

func (r *Resolver) ProductReview() gqlgen.ProductReviewResolver {
	return r.ReviewResolver.ProductReviewResolver
}

func (r *Resolver) ServiceReview() gqlgen.ServiceReviewResolver {
	return r.ReviewResolver.ServiceReviewResolver
}

func (r *Resolver) AppointmentReview() gqlgen.AppointmentReviewResolver {
	return r.ReviewResolver.AppointmentReviewResolver
}

func (r *Resolver) Favorite() gqlgen.FavoriteResolver {
	return r.FavoriteResolver
}

func (r *Resolver) News() gqlgen.NewsResolver {
	return r.NewsResolver
}

func (r *Resolver) EmailTemplate() gqlgen.EmailTemplateResolver {
	return r.EmailTemplateResolver
}

func (r *Resolver) ProductServiceAttribute() gqlgen.ProductServiceAttributeResolver {
	return r.ProductServiceAttributeResolver
}

var ProviderSet = wire.NewSet(
	New,
)
