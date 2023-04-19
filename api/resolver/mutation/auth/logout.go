package auth

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (a *Auth) Logout(ctx context.Context) (*gqlgen.LogoutPayload, error) {
	token, err := sessctx.Token(ctx)

	if err != nil {
		return nil, err
	}

	_, err = a.Prisma.DeleteSession(prisma.SessionWhereUniqueInput{
		Token: &token,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.LogoutPayload{
		Session: nil,
	}, nil
}
