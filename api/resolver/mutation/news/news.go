package news

import (
	"github.com/steebchen/keskin-api/prisma"
)

type NewsMutation struct {
	Prisma *prisma.Client
}

var allowedTypes = []prisma.UserType{
	prisma.UserTypeManager,
}

func New(client *prisma.Client) *NewsMutation {
	return &NewsMutation{
		Prisma: client,
	}
}
