package employee

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Employee) FirstName(ctx context.Context, obj *prisma.Employee) (*string, error) {
	if obj.Deleted {
		return nil, nil
	}

	return &obj.FirstName, nil
}
