package branch

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Branch) Company(ctx context.Context, obj *prisma.Branch) (*prisma.Company, error) {
	return r.Prisma.Branch(prisma.BranchWhereUniqueInput{
		ID: &obj.ID,
	}).Company().Exec(ctx)
}
