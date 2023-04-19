package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) AfterImage(ctx context.Context, obj *prisma.Appointment) (*gqlgen.Image, error) {
	return picture.FromID(obj.AfterImage), nil
}
