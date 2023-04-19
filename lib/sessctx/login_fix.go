package sessctx

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

type LoginFixData struct {
	Apply        bool
	SessionToken string
	User         *prisma.User
}

func NewLoginFix() *LoginFixData {
	return &LoginFixData{
		false,
		"",
		nil,
	}
}

func InitLoginFix(ctx context.Context) context.Context {
	return context.WithValue(ctx, "loginFix", NewLoginFix())
}

func LoginFix(ctx context.Context) *LoginFixData {
	return ctx.Value("loginFix").(*LoginFixData)
}
