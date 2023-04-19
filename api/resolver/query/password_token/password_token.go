package password_token

import (
	"context"
	"time"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func New(client *prisma.Client) *PasswordTokenQuery {
	return &PasswordTokenQuery{
		Prisma: client,
	}
}

type PasswordTokenQuery struct {
	Prisma *prisma.Client
}

func QueryValidPasswordTokens(prismaClient *prisma.Client, ctx context.Context, token string) ([]prisma.PasswordToken, error) {
	invalidBefore := time.Now().Add(-time.Duration(24) * time.Hour)
	invalidBeforeDate := invalidBefore.Format(time.RFC3339)

	deleted := false

	return prismaClient.PasswordTokens(&prisma.PasswordTokensParams{
		Where: &prisma.PasswordTokenWhereInput{
			Token:        &token,
			CreatedAtGte: &invalidBeforeDate,
			User: &prisma.UserWhereInput{
				Deleted: &deleted,
			},
		},
	}).Exec(ctx)
}

func (r *PasswordTokenQuery) IsValidPasswordToken(
	ctx context.Context,
	token string,
) (*gqlgen.IsValidPasswordTokenPayload, error) {
	passwordTokens, err := QueryValidPasswordTokens(r.Prisma, ctx, token)

	if err != nil && err != prisma.ErrNoResult {
		return &gqlgen.IsValidPasswordTokenPayload{
			Valid: false,
		}, err
	}

	return &gqlgen.IsValidPasswordTokenPayload{
		Valid: len(passwordTokens) > 0,
	}, nil
}
