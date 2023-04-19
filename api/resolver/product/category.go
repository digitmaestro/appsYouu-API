package product

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Product) Category(ctx context.Context, obj *prisma.Product) (*string, error) {
	if obj.Deleted {
		return nil, nil
	}

	return obj.Category, nil
}
