package auth

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
)

func (a *Auth) RequestActivationLink(
	ctx context.Context,
	input gqlgen.RequestActivationLinkInput,
) (*gqlgen.RequestActivationLinkPayload, error) {
	RequestActivationLink(a.Prisma, ctx, input.Email, input.Company)

	return &gqlgen.RequestActivationLinkPayload{
		Status: "OK",
	}, nil
}
