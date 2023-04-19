package appointment_review

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/prisma"
)

type AppointmentReview struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *AppointmentReview {
	return &AppointmentReview{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
