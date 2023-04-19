package customer

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Customer) FirstName(ctx context.Context, obj *prisma.Customer) (*string, error) {
	if obj.Deleted {
		return nil, nil
	}

	return &obj.FirstName, nil
}
