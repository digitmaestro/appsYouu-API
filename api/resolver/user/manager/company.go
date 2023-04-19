package manager

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Manager) Company(ctx context.Context, obj *prisma.Manager) (*prisma.Company, error) {
	if obj.Deleted {
		return nil, nil
	}

	return r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: &obj.ID,
	}).Company().Exec(ctx)
}
