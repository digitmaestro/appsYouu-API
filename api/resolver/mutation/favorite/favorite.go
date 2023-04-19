package favorite

import (
	"github.com/steebchen/keskin-api/prisma"
)

type FavoriteMutation struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *FavoriteMutation {
	return &FavoriteMutation{
		Prisma: client,
	}
}
