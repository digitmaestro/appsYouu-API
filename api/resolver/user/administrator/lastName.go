package administrator

import (
	"context"

	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Administrator) LastName(ctx context.Context, obj *prisma.Administrator) (*string, error) {
	if obj.Deleted {
		deleted := i18n.Language(ctx)["DELETED_USER"]
		return &deleted, nil
	}

	return &obj.LastName, nil
}
