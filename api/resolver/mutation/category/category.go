package category

import "github.com/steebchen/keskin-api/prisma"

type CategoryMutation struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *CategoryMutation {
	return &CategoryMutation{
		Prisma: client,
	}
}
