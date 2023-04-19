package manager

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/phonenumber"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Manager) PhoneNumber(ctx context.Context, obj *prisma.Manager) (*gqlgen.PhoneNumber, error) {
	viewer, err := sessctx.User(ctx)

	if obj.Deleted || err == sessctx.UserNotLoggedInError {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if viewer.Type == prisma.UserTypeEmployee || viewer.Type == prisma.UserTypeManager || viewer.Type == prisma.UserTypeAdministrator || viewer.ID == obj.ID {
		return phonenumber.Convert(obj.PhoneNumber), nil
	} else {
		return nil, nil
	}
}
