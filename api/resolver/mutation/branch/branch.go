package branch

import (
	"github.com/steebchen/keskin-api/prisma"
)

type BranchMutation struct {
	Prisma *prisma.Client
}

var allowedTypes = []prisma.UserType{
	prisma.UserTypeManager,
}

func New(client *prisma.Client) *BranchMutation {
	return &BranchMutation{
		Prisma: client,
	}
}
