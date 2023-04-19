package service

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

type Service struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *Service {
	return &Service{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
