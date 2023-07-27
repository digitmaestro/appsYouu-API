package auth

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/validator"
)

func (a *Auth) RequestPasswordReset(
	ctx context.Context,
	input gqlgen.RequestPasswordResetInput,
) (*gqlgen.RequestPasswordResetPayload, error) {

	err := validator.Email(input.Email)
	if err != nil {
		return nil, gqlerrors.NewValidationError(err.Error(), "InvalidEmail")
	}

	RequestPasswordReset(a.Prisma, ctx, input.Email, input.Company)

	return &gqlgen.RequestPasswordResetPayload{
		Status: "OK",
	}, nil
}
