package service_review

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

type ServiceReview struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *ServiceReview {
	return &ServiceReview{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
