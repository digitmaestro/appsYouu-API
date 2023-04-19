package branch

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *BranchMutation) DeleteBranch(
	ctx context.Context,
	input gqlgen.DeleteBranchInput,
	language *string,
) (*gqlgen.DeleteBranchPayload, error) {
	company, err := r.Prisma.Branch(prisma.BranchWhereUniqueInput{
		ID: &input.ID,
	}).Company().Exec(ctx)

	if err != nil {
		return nil, err
	}

	if err := permissions.CanAccessCompany(ctx, company.ID, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	branch, err := r.Prisma.DeleteBranch(prisma.BranchWhereUniqueInput{
		ID: &input.ID,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.DeleteBranchPayload{
		Branch: branch,
	}, nil
}
