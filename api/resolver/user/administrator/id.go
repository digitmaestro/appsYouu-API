package administrator

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Administrator) ID(ctx context.Context, obj *prisma.Administrator) (string, error) {
	if obj.Deleted {
		return "", nil
	}

	return obj.ID, nil
}
