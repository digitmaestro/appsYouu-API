package auth

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/validator"
)

func (a *Auth) RequestActivationLink(
	ctx context.Context,
	input gqlgen.RequestActivationLinkInput,
) (*gqlgen.RequestActivationLinkPayload, error) {

	err := validator.Email(input.Email)
	if err != nil {
		return nil, gqlerrors.NewValidationError(err.Error(), "InvalidEmail")
	}

	RequestActivationLink(a.Prisma, ctx, input.Email, input.Company)

	return &gqlgen.RequestActivationLinkPayload{
		Status: "OK",
	}, nil
}
