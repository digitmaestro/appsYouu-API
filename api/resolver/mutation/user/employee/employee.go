package employee

import (
	"github.com/steebchen/keskin-api/prisma"
)

type EmployeeMutation struct {
	Prisma *prisma.Client
}

var allowedTypes = []prisma.UserType{
	prisma.UserTypeManager,
}

func New(client *prisma.Client) *EmployeeMutation {
	return &EmployeeMutation{
		Prisma: client,
	}
}
