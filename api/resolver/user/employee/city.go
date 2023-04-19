package employee

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Employee) City(ctx context.Context, obj *prisma.Employee) (*string, error) {
	viewer, err := sessctx.User(ctx)

	if obj.Deleted || err != nil {
		return nil, err
	}

	if viewer.Type == prisma.UserTypeEmployee || viewer.Type == prisma.UserTypeManager || viewer.Type == prisma.UserTypeAdministrator || viewer.ID == obj.ID {
		return obj.City, nil
	} else {
		return nil, nil
	}
}
