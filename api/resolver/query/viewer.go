package query

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Query) Viewer(ctx context.Context) (prisma.IUser, error) {
	user, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	return user.Convert(), err
}
