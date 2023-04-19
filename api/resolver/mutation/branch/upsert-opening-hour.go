package branch

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *BranchMutation) UpsertBranchOpeningHour(
	ctx context.Context,
	input gqlgen.UpsertBranchOpeningHourInput,
) (*gqlgen.UpsertBranchOpeningHourPayload, error) {
	if err := permissions.CanAccessBranch(ctx, input.Branch, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	// only allow one unique day per branch
	// cut id because prisma limits id length
	id := prisma.Str(input.Branch[10:] + "-" + string(input.Patch.Day))

	connect := &prisma.BranchWhereUniqueInput{
		ID: &input.Branch,
	}

	create := prisma.BranchOpeningHourCreateInput{
		ID: id,
		Branch: prisma.BranchCreateOneWithoutOpeningHoursInput{
			Connect: connect,
		},
		Day:    input.Patch.Day,
		Closed: &input.Patch.Closed,
		Break:  &input.Patch.Break,
	}

	update := prisma.BranchOpeningHourUpdateInput{
		Branch: &prisma.BranchUpdateOneRequiredWithoutOpeningHoursInput{
			Connect: connect,
		},
		Day:    &input.Patch.Day,
		Closed: &input.Patch.Closed,
		Break:  &input.Patch.Break,
	}

	if !input.Patch.Closed {
		if input.Patch.Forenoon != nil {
			create.StartForenoon = prisma.TimeString(input.Patch.Forenoon.Start)
			create.EndForenoon = prisma.TimeString(input.Patch.Forenoon.End)
			update.StartForenoon = prisma.TimeString(input.Patch.Forenoon.Start)
			update.EndForenoon = prisma.TimeString(input.Patch.Forenoon.End)
		}

		if input.Patch.Break && input.Patch.Afternoon != nil {
			create.StartAfternoon = prisma.TimeString(input.Patch.Afternoon.Start)
			create.EndAfternoon = prisma.TimeString(input.Patch.Afternoon.End)
			update.StartAfternoon = prisma.TimeString(input.Patch.Afternoon.Start)
			update.EndAfternoon = prisma.TimeString(input.Patch.Afternoon.End)
		}
	}

	_, err := r.Prisma.UpsertBranchOpeningHour(prisma.BranchOpeningHourUpsertParams{
		Where: prisma.BranchOpeningHourWhereUniqueInput{
			ID: id,
		},
		Create: create,
		Update: update,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	branch, err := r.Prisma.Branch(prisma.BranchWhereUniqueInput{
		ID: &input.Branch,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpsertBranchOpeningHourPayload{
		Branch: branch,
	}, nil
}
