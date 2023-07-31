package sub_category

import (
	"github.com/steebchen/keskin-api/prisma"
)

type SubCategoryMutation struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *SubCategoryMutation {
	return &SubCategoryMutation{
		Prisma: client,
	}
}
