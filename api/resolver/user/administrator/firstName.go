package administrator

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Administrator) FirstName(ctx context.Context, obj *prisma.Administrator) (*string, error) {
	if obj.Deleted {
		return nil, nil
	}

	return &obj.FirstName, nil
}
