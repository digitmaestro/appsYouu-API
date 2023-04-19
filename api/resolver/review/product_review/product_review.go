package product_review

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

type ProductReview struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *ProductReview {
	return &ProductReview{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
