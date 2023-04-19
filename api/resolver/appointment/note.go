package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) Note(ctx context.Context, obj *prisma.Appointment) (*string, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	switch viewer.Type {
	case prisma.UserTypeEmployee, prisma.UserTypeManager, prisma.UserTypeAdministrator:
		return obj.Note, nil
	default:
		return nil, err
	}
}
