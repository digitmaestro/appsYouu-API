package customer

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Customer) Note(ctx context.Context, obj *prisma.Customer) (*string, error) {
	viewer, err := sessctx.User(ctx)

	if obj.Deleted || err != nil {
		return nil, err
	}

	if viewer.Type == prisma.UserTypeEmployee || viewer.Type == prisma.UserTypeManager || viewer.Type == prisma.UserTypeAdministrator {
		return obj.Note, nil
	} else {
		return nil, nil
	}
}
