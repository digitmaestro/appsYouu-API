package administrator

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Administrator) ZipCode(ctx context.Context, obj *prisma.Administrator) (*string, error) {
	viewer, err := sessctx.User(ctx)

	if obj.Deleted || err != nil {
		return nil, err
	}

	if viewer.Type == prisma.UserTypeEmployee || viewer.Type == prisma.UserTypeManager || viewer.Type == prisma.UserTypeAdministrator || viewer.ID == obj.ID {
		return obj.ZipCode, nil
	} else {
		return nil, nil
	}
}
