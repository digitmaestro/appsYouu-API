package appointment

import (
	"context"
	"time"

	"github.com/steebchen/keskin-api/prisma"
)

var utc, _ = time.LoadLocation("UTC")

func (r *Appointment) Duration(ctx context.Context, obj *prisma.Appointment) (*int, error) {
	duration := prisma.TimeDate(utc, obj.End).Sub(prisma.TimeDate(utc, obj.Start))
	m := int(duration.Minutes())
	return &m, nil
}
