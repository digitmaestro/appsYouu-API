package auth

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/query/password_token"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/auth"
	"github.com/steebchen/keskin-api/prisma"
)

func (a *Auth) ResetPassword(
	ctx context.Context,
	input gqlgen.ResetPasswordInput,
) (*gqlgen.ResetPasswordPayload, error) {
	passwordTokens, err := password_token.QueryValidPasswordTokens(a.Prisma, ctx, input.Token)

	if err != nil || len(passwordTokens) == 0 {
		return nil, gqlerrors.NewPermissionError("invalid password reset token")
	}

	user, err := a.Prisma.PasswordToken(prisma.PasswordTokenWhereUniqueInput{
		ID: &passwordTokens[0].ID,
	}).User().Exec(ctx)

	if err != nil || user.Deleted {
		return nil, gqlerrors.NewPermissionError("invalid password reset token")
	}

	_, err = a.Prisma.DeletePasswordToken(prisma.PasswordTokenWhereUniqueInput{
		ID: &passwordTokens[0].ID,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	passwordHash := auth.HashPassword(input.Password)

	_, err = a.Prisma.UpdateUser(prisma.UserUpdateParams{
		Where: prisma.UserWhereUniqueInput{
			ID: &user.ID,
		},
		Data: prisma.UserUpdateInput{
			PasswordHash: &passwordHash,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.ResetPasswordPayload{
		Status: "OK",
	}, nil
}
