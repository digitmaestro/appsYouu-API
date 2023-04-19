package viewer

import (
	"github.com/steebchen/keskin-api/prisma"
)

type ViewerMutation struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *ViewerMutation {
	return &ViewerMutation{
		Prisma: client,
	}
}
