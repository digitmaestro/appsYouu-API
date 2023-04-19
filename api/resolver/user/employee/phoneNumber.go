package employee

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/phonenumber"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Employee) PhoneNumber(ctx context.Context, obj *prisma.Employee) (*gqlgen.PhoneNumber, error) {
	viewer, err := sessctx.User(ctx)

	if obj.Deleted || err != nil {
		return nil, err
	}

	if viewer.Type == prisma.UserTypeEmployee || viewer.Type == prisma.UserTypeManager || viewer.Type == prisma.UserTypeAdministrator || viewer.ID == obj.ID {
		return phonenumber.Convert(obj.PhoneNumber), nil
	} else {
		return nil, nil
	}
}
