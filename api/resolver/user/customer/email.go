package customer

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Customer) Email(ctx context.Context, obj *prisma.Customer) (string, error) {
	viewer, err := sessctx.User(ctx)

	if obj.Deleted || err != nil {
		return "", err
	}

	if viewer.Type == prisma.UserTypeEmployee || viewer.Type == prisma.UserTypeManager || viewer.Type == prisma.UserTypeAdministrator || viewer.ID == obj.ID {
		return obj.Email, nil
	} else {
		return "", nil
	}
}
