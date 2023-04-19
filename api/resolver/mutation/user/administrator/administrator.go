package administrator

import (
	"github.com/steebchen/keskin-api/prisma"
)

type AdministratorMutation struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *AdministratorMutation {
	return &AdministratorMutation{
		Prisma: client,
	}
}
