package auth

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
)

func (a *Auth) RequestPasswordReset(
	ctx context.Context,
	input gqlgen.RequestPasswordResetInput,
) (*gqlgen.RequestPasswordResetPayload, error) {
	RequestPasswordReset(a.Prisma, ctx, input.Email, input.Company)

	return &gqlgen.RequestPasswordResetPayload{
		Status: "OK",
	}, nil
}
