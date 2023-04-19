package email_template

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

type EmailTemplate struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *EmailTemplate {
	return &EmailTemplate{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
