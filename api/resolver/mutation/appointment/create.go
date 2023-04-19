package appointment

import (
	"context"
	"time"

	"github.com/steebchen/keskin-api/api/resolver/appointment"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

type CreateAppointmentInput struct {
	Client          *prisma.Client
	Context         context.Context
	EmployeeID      *string
	CustomerID      string
	Branch          string
	Desc            string
	Start           time.Time
	ProductRequests []*gqlgen.ConnectAppointmentProduct
	ServiceRequests []*gqlgen.ConnectAppointmentService
	DefaultStatus   prisma.AppointmentStatus
}

func CreateAppointment(input CreateAppointmentInput) (*prisma.Appointment, error) {
	if len(input.ServiceRequests) == 0 {
		return nil, gqlerrors.NewInternalError("services require at least one service", "OneServiceRequired")
	}

	if input.Start.Before(time.Now()) {
		return nil, gqlerrors.NewValidationError("Start date can not be in the past", "InvalidTime")
	}

	data, err := appointment.Plan(appointment.PlanInput{
		Client:          input.Client,
		Context:         input.Context,
		ProductRequests: input.ProductRequests,
		ServiceRequests: input.ServiceRequests,
		Start:           input.Start,
	})

	if input.EmployeeID == nil {
		employee, err := FindAvailableEmployee(input.Client, input.Context, input.Branch, &data.Start, &data.End)

		if err != nil {
			return nil, err
		}

		input.EmployeeID = &employee.ID
	}

	if err != nil {
		return nil, err
	}

	customer, err := input.Client.User(prisma.UserWhereUniqueInput{
		ID: &input.CustomerID,
	}).Exec(input.Context)

	if err != nil {
		return nil, err
	}

	if customer.Deleted {
		return nil, gqlerrors.NewPermissionError("Customer is deleted")
	}

	apt, err := input.Client.CreateAppointment(prisma.AppointmentCreateInput{
		Start:  data.Start,
		End:    data.End,
		Desc:   *i18n.CreateLocalizedString(input.Context, &input.Desc),
		Status: input.DefaultStatus,

		Price: data.Price,

		Products: &prisma.AppointmentProductLinkCreateManyWithoutAppointmentInput{
			Create: data.CreateProducts,
		},

		Services: &prisma.AppointmentServiceLinkCreateManyWithoutAppointmentInput{
			Create: data.CreateServices,
		},

		Customer: prisma.UserCreateOneInput{
			Connect: &prisma.UserWhereUniqueInput{
				ID: &input.CustomerID,
			},
		},
		Employee: prisma.UserCreateOneInput{
			Connect: &prisma.UserWhereUniqueInput{
				ID: input.EmployeeID,
			},
		},
		Branch: prisma.BranchCreateOneInput{
			Connect: &prisma.BranchWhereUniqueInput{
				ID: &input.Branch,
			},
		},
	}).Exec(input.Context)

	if err != nil {
		return nil, err
	}

	return apt, nil
}
