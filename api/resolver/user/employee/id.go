package employee

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Employee) ID(ctx context.Context, obj *prisma.Employee) (string, error) {
	if obj.Deleted {
		return "", nil
	}

	return obj.ID, nil
}
