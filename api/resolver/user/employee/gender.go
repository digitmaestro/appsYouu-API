package employee

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Employee) Gender(ctx context.Context, obj *prisma.Employee) (*prisma.Gender, error) {
	if obj.Deleted {
		return nil, nil
	}

	return &obj.Gender, nil
}
