package customer

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Customer) ID(ctx context.Context, obj *prisma.Customer) (string, error) {
	if obj.Deleted {
		return "", nil
	}

	return obj.ID, nil
}
