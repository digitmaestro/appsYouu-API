package hours

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func UpsertWorkingHours(prismaClient *prisma.Client, ctx context.Context, userId string, workingHours []*gqlgen.WorkingHoursInput) error {
	for _, workingHoursItem := range workingHours {
		existingWorkingHours, err := prismaClient.WorkingHourses(&prisma.WorkingHoursesParams{
			Where: &prisma.WorkingHoursWhereInput{
				Day: &workingHoursItem.Day,
				User: &prisma.UserWhereInput{
					ID: &userId,
				},
			},
		}).Exec(ctx)

		if err != nil {
			return err
		}

		if len(existingWorkingHours) == 0 {
			create := prisma.WorkingHoursCreateInput{
				Day:        workingHoursItem.Day,
				NotWorking: &workingHoursItem.NotWorking,
				Break:      &workingHoursItem.Break,
				User: prisma.UserCreateOneWithoutWorkingHoursInput{
					Connect: &prisma.UserWhereUniqueInput{
						ID: &userId,
					},
				},
			}

			if workingHoursItem.Forenoon != nil {
				create.StartForenoon = prisma.TimeString(workingHoursItem.Forenoon.Start)
				create.EndForenoon = prisma.TimeString(workingHoursItem.Forenoon.End)
			}

			if workingHoursItem.Afternoon != nil {
				create.StartAfternoon = prisma.TimeString(workingHoursItem.Afternoon.Start)
				create.EndAfternoon = prisma.TimeString(workingHoursItem.Afternoon.End)
			}

			_, err = prismaClient.CreateWorkingHours(create).Exec(ctx)
		} else {
			update := prisma.WorkingHoursUpdateInput{
				NotWorking: &workingHoursItem.NotWorking,
				Break:      &workingHoursItem.Break,
			}

			if workingHoursItem.Forenoon != nil {
				update.StartForenoon = prisma.TimeString(workingHoursItem.Forenoon.Start)
				update.EndForenoon = prisma.TimeString(workingHoursItem.Forenoon.End)
			}

			if workingHoursItem.Afternoon != nil {
				update.StartAfternoon = prisma.TimeString(workingHoursItem.Afternoon.Start)
				update.EndAfternoon = prisma.TimeString(workingHoursItem.Afternoon.End)
			}

			_, err = prismaClient.UpdateWorkingHours(prisma.WorkingHoursUpdateParams{
				Where: prisma.WorkingHoursWhereUniqueInput{
					ID: &existingWorkingHours[0].ID,
				},
				Data: update,
			}).Exec(ctx)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
