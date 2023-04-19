package appointment

import (
	"context"
	"errors"
	"time"

	"github.com/steebchen/keskin-api/lib/hours"
	"github.com/steebchen/keskin-api/prisma"
)

var utc, _ = time.LoadLocation("UTC")

func FindAvailableEmployee(client *prisma.Client, ctx context.Context, branch string, start *string, end *string) (*prisma.Employee, error) {
	userTypeEmployee := prisma.UserTypeEmployee
	userTypeManager := prisma.UserTypeManager
	notWorking := false
	deleted := false
	dayOfWeek := hours.WeekdayToDayOfWeek(prisma.TimeDate(utc, *start).Weekday())

	workingHours, err := client.WorkingHourses(&prisma.WorkingHoursesParams{
		Where: &prisma.WorkingHoursWhereInput{
			User: &prisma.UserWhereInput{
				Deleted: &deleted,
				Or: []prisma.UserWhereInput{
					{
						Type: &userTypeEmployee,
						Branch: &prisma.BranchWhereInput{
							ID: &branch,
						},
					},
					{
						Type: &userTypeManager,
						Company: &prisma.CompanyWhereInput{
							BranchesSome: &prisma.BranchWhereInput{
								ID: &branch,
							},
						},
					},
				},
			},
			Day:        &dayOfWeek,
			NotWorking: &notWorking,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(workingHours) == 0 {
		return nil, errors.New("could not find an employee")
	}

	workingEmployeesMap := make(map[string]bool)

	startToday := hours.TodaySetTime(start)
	endToday := hours.TodaySetTime(end)

	for _, workingHoursItem := range workingHours {
		isWithinWorkingHours := false

		forenoonStart := hours.TodaySetTime(workingHoursItem.StartForenoon)
		forenoonEnd := hours.TodaySetTime(workingHoursItem.EndForenoon)

		if (forenoonStart.Before(startToday) ||
			forenoonStart.Equal(startToday)) &&
			(forenoonEnd.After(endToday) ||
				forenoonEnd.Equal(endToday)) {
			isWithinWorkingHours = true
		}

		if !isWithinWorkingHours && workingHoursItem.Break {
			afternoonStart := hours.TodaySetTime(workingHoursItem.StartAfternoon)
			afternoonEnd := hours.TodaySetTime(workingHoursItem.EndAfternoon)

			if (afternoonStart.Before(startToday) ||
				afternoonStart.Equal(startToday)) &&
				(afternoonEnd.After(endToday) ||
					afternoonEnd.Equal(endToday)) {
				isWithinWorkingHours = true
			}
		}

		if isWithinWorkingHours {
			user, err := client.WorkingHours(prisma.WorkingHoursWhereUniqueInput{
				ID: &workingHoursItem.ID,
			}).User().Exec(ctx)

			if err == nil {
				workingEmployeesMap[user.ID] = true
			}
		}
	}

	if len(workingEmployeesMap) == 0 {
		return nil, errors.New("could not find an employee")
	}

	appointments, err := client.Appointments(&prisma.AppointmentsParams{
		Where: &prisma.AppointmentWhereInput{
			Branch: &prisma.BranchWhereInput{
				ID: &branch,
			},
			StatusIn: []prisma.AppointmentStatus{
				prisma.AppointmentStatusApproved,
				prisma.AppointmentStatusRequested,
			},
			Or: []prisma.AppointmentWhereInput{
				{
					And: []prisma.AppointmentWhereInput{
						{
							StartLte: start,
						},
						{
							EndGte: start,
						},
					},
				},
				{
					And: []prisma.AppointmentWhereInput{
						{
							StartLte: end,
						},
						{
							EndGte: end,
						},
					},
				},
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	for _, appointment := range appointments {
		user, err := client.Appointment(prisma.AppointmentWhereUniqueInput{
			ID: &appointment.ID,
		}).Employee().Exec(ctx)

		if err == nil {
			workingEmployeesMap[user.ID] = false
		}
	}

	workingEmployees := []string{}

	for userId, available := range workingEmployeesMap {
		if available {
			workingEmployees = append(workingEmployees, userId)
		}
	}

	if len(workingEmployees) == 0 {
		return nil, errors.New("could not find an employee")
	}

	employees, err := client.Users(&prisma.UsersParams{
		First: prisma.Int32(1),
		Where: &prisma.UserWhereInput{
			Type: &userTypeEmployee,
			Branch: &prisma.BranchWhereInput{
				ID: &branch,
			},
			IDIn:    workingEmployees,
			Deleted: &deleted,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(employees) == 0 {
		employees, err = client.Users(&prisma.UsersParams{
			First: prisma.Int32(1),
			Where: &prisma.UserWhereInput{
				Type: &userTypeManager,
				Company: &prisma.CompanyWhereInput{
					BranchesSome: &prisma.BranchWhereInput{
						ID: &branch,
					},
				},
				IDIn:    workingEmployees,
				Deleted: &deleted,
			},
		}).Exec(ctx)
	}

	if err != nil {
		return nil, err
	}

	if len(employees) == 0 {
		return nil, errors.New("could not find an employee")
	}

	return &prisma.Employee{
		User: &employees[0],
	}, nil
}
