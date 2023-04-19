package product

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Product) ID(ctx context.Context, obj *prisma.Product) (string, error) {
	if obj.Deleted {
		return "", nil
	}

	return obj.ID, nil
}
