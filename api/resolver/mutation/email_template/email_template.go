package email_template

import (
	"github.com/steebchen/keskin-api/prisma"
)

var allowedTypes = []prisma.UserType{
	prisma.UserTypeManager,
	prisma.UserTypeEmployee,
}

type EmailTemplateMutation struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *EmailTemplateMutation {
	return &EmailTemplateMutation{
		Prisma: client,
	}
}
