package timeslots

import (
	"context"
	"time"

	"github.com/steebchen/keskin-api/api/resolver/appointment"
	"github.com/steebchen/keskin-api/api/resolver/price"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/prisma"
)

func New(client *prisma.Client) *TimeslotsQuery {
	return &TimeslotsQuery{
		Prisma: client,
	}
}

type TimeslotsQuery struct {
	Prisma *prisma.Client
}

func startOfDay(t time.Time) time.Time {
	d := time.Date(
		t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0,
		utc,
	)
	return d
}

func endOfDay(t time.Time) time.Time {
	d := time.Date(
		t.Year(), t.Month(), t.Day(),
		23, 59, 59, 999,
		utc,
	)
	return d
}

func (r *TimeslotsQuery) AppointmentTimeslots(
	ctx context.Context,
	input gqlgen.TimeslotInput,
) (*gqlgen.Timeslots, error) {
	var start time.Time
	if input.Start == nil {
		start = time.Now()
	} else {
		if input.Start != nil &&
			input.Start.Before(time.Now().Truncate(24*time.Hour)) {
			return nil, gqlerrors.NewValidationError("Start date can not be in the past", "InvalidTime")
		}

		start = *input.Start
	}

	_, err := r.Prisma.Branch(prisma.BranchWhereUniqueInput{
		ID: &input.Branch,
	}).Exec(ctx)

	if err != nil {
		return nil, gqlerrors.NewFormatNodeError(err, input.Branch)
	}

	deleted := false

	if input.Employee != nil {
		staff, err := r.Prisma.User(prisma.UserWhereUniqueInput{
			ID: input.Employee,
		}).Exec(ctx)

		if err != nil {
			return nil, gqlerrors.NewFormatNodeError(err, *input.Employee)
		}

		if staff.Type != prisma.UserTypeManager && staff.Type != prisma.UserTypeEmployee {
			return nil, gqlerrors.NewValidationError("staff type is invalid", "InvalidUserType")
		}

		if staff.Deleted {
			return nil, gqlerrors.NewPermissionError("Staff is deleted")
		}
	}

	openingHours, err := r.Prisma.BranchOpeningHours(&prisma.BranchOpeningHoursParams{
		Where: &prisma.BranchOpeningHourWhereInput{
			Branch: &prisma.BranchWhereInput{
				ID: &input.Branch,
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	workingHoursWhere := &prisma.WorkingHoursWhereInput{}
	if input.Employee != nil {
		workingHoursWhere.User = &prisma.UserWhereInput{
			ID:      input.Employee,
			Deleted: &deleted,
		}
	} else {
		userTypeEmployee := prisma.UserTypeEmployee
		userTypeManager := prisma.UserTypeManager

		workingHoursWhere.User = &prisma.UserWhereInput{
			Deleted: &deleted,
			Or: []prisma.UserWhereInput{
				{
					Type: &userTypeEmployee,
					Branch: &prisma.BranchWhereInput{
						ID: &input.Branch,
					},
				},
				{
					Type: &userTypeManager,
					Company: &prisma.CompanyWhereInput{
						BranchesSome: &prisma.BranchWhereInput{
							ID: &input.Branch,
						},
					},
				},
			},
		}
	}

	workingHours, err := r.Prisma.WorkingHourses(&prisma.WorkingHoursesParams{
		Where: workingHoursWhere,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	appointmentsWhere := prisma.AppointmentWhereInput{
		Branch: &prisma.BranchWhereInput{
			ID: &input.Branch,
		},
		StartGte: prisma.TimeString(startOfDay(start)),
		EndLte:   prisma.TimeString(endOfDay(start)),
		StatusIn: []prisma.AppointmentStatus{
			prisma.AppointmentStatusApproved,
			prisma.AppointmentStatusRequested,
		},
	}

	if input.Employee != nil {
		appointmentsWhere.Employee = &prisma.UserWhereInput{
			ID: input.Employee,
		}
	}

	appointments, err := r.Prisma.Appointments(&prisma.AppointmentsParams{
		Where: &appointmentsWhere,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	data, err := appointment.Plan(appointment.PlanInput{
		Client:          r.Prisma,
		Context:         ctx,
		ProductRequests: input.Products,
		ServiceRequests: input.Services,
		Start:           start,
	})

	if err != nil {
		return nil, err
	}

	openingHourToday := findOpeningHour(input.Start, openingHours)
	workingHoursToday := findWorkingHours(input.Start, workingHours)

	ranges := timeslotRanges(
		start,
		data.Duration,
		openingHourToday,
		calculateWorkingHoursSlices(
			r.Prisma,
			start,
			appointments,
			workingHoursToday,
		),
	)

	nodes := []*gqlgen.TimeslotRange{}

	for _, item := range ranges {
		clone := item
		nodes = append(nodes, &clone)
	}

	return &gqlgen.Timeslots{
		Price:    price.Convert(ctx, data.Price),
		Duration: int(data.Duration.Minutes()),
		Ranges:   nodes,
	}, nil
}
