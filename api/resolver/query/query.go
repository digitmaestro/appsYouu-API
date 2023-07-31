package query

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/api/resolver/query/appointment"
	"github.com/steebchen/keskin-api/api/resolver/query/branch"
	"github.com/steebchen/keskin-api/api/resolver/query/category"
	"github.com/steebchen/keskin-api/api/resolver/query/company"
	"github.com/steebchen/keskin-api/api/resolver/query/customer"
	"github.com/steebchen/keskin-api/api/resolver/query/email_template"
	"github.com/steebchen/keskin-api/api/resolver/query/employee"
	"github.com/steebchen/keskin-api/api/resolver/query/favorite"
	"github.com/steebchen/keskin-api/api/resolver/query/order_history"
	"github.com/steebchen/keskin-api/api/resolver/query/password_token"
	"github.com/steebchen/keskin-api/api/resolver/query/product"
	"github.com/steebchen/keskin-api/api/resolver/query/review"
	"github.com/steebchen/keskin-api/api/resolver/query/service"
	"github.com/steebchen/keskin-api/api/resolver/query/staff"
	"github.com/steebchen/keskin-api/api/resolver/query/sub_category"
	"github.com/steebchen/keskin-api/api/resolver/query/timeslots"
	"github.com/steebchen/keskin-api/prisma"
)

type Query struct {
	Prisma *prisma.Client
	*appointment.AppointmentQuery
	*branch.BranchQuery
	*customer.CustomerQuery
	*category.CategoryQuery
	*sub_category.SubCategoryQuery
	*employee.EmployeeQuery
	*product.ProductQuery
	*service.ServiceQuery
	*staff.StaffQuery
	*timeslots.TimeslotsQuery
	*review.ReviewQuery
	*order_history.OrderHistoryQuery
	*favorite.FavoriteQuery
	*email_template.EmailTemplateQuery
	*company.CompanyQuery
	*password_token.PasswordTokenQuery
}

func New(
	client *prisma.Client,
	appointment *appointment.AppointmentQuery,
	branch *branch.BranchQuery,
	customer *customer.CustomerQuery,
	category *category.CategoryQuery,
	sub_category *sub_category.SubCategoryQuery,
	employee *employee.EmployeeQuery,
	product *product.ProductQuery,
	service *service.ServiceQuery,
	staff *staff.StaffQuery,
	timeslots *timeslots.TimeslotsQuery,
	review *review.ReviewQuery,
	orderHistory *order_history.OrderHistoryQuery,
	favorite *favorite.FavoriteQuery,
	emailTemplate *email_template.EmailTemplateQuery,
	company *company.CompanyQuery,
	passwordToken *password_token.PasswordTokenQuery,
) *Query {
	return &Query{
		client,
		appointment,
		branch,
		customer,
		category,
		sub_category,
		employee,
		product,
		service,
		staff,
		timeslots,
		review,
		orderHistory,
		favorite,
		emailTemplate,
		company,
		passwordToken,
	}
}

var ProviderSet = wire.NewSet(
	New,
	appointment.New,
	branch.New,
	customer.New,
	employee.New,
	category.New,
	sub_category.New,
	product.New,
	service.New,
	staff.New,
	timeslots.New,
	review.New,
	order_history.New,
	favorite.New,
	email_template.New,
	company.New,
	password_token.New,
)
