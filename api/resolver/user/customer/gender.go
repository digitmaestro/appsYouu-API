package customer

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Customer) Gender(ctx context.Context, obj *prisma.Customer) (*prisma.Gender, error) {
	if obj.Deleted {
		return nil, nil
	}

	return &obj.Gender, nil
}
