package company

import (
	"github.com/steebchen/keskin-api/prisma"
)

type CompanyMutation struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *CompanyMutation {
	return &CompanyMutation{
		Prisma: client,
	}
}
