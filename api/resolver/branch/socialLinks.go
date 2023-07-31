package branch

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Branch) FacebookLink(ctx context.Context, obj *prisma.Branch) (*string, error) {
	branch, _ := r.Prisma.Branch(prisma.BranchWhereUniqueInput{
		ID: &obj.ID,
	}).Exec(ctx)

	return branch.FacebookLink, nil
}

func (r *Branch) TiktokLink(ctx context.Context, obj *prisma.Branch) (*string, error) {
	branch, _ := r.Prisma.Branch(prisma.BranchWhereUniqueInput{
		ID: &obj.ID,
	}).Exec(ctx)

	return branch.TiktokLink, nil
}

func (r *Branch) InstagramLink(ctx context.Context, obj *prisma.Branch) (*string, error) {
	branch, _ := r.Prisma.Branch(prisma.BranchWhereUniqueInput{
		ID: &obj.ID,
	}).Exec(ctx)

	return branch.InstagramLink, nil
}
