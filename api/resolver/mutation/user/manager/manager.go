package manager

import (
	"github.com/steebchen/keskin-api/prisma"
)

var allowedTypes = []prisma.UserType{
	prisma.UserTypeManager,
}

type ManagerMutation struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *ManagerMutation {
	return &ManagerMutation{
		Prisma: client,
	}
}
