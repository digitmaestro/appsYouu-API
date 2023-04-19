package administrator

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Administrator) Gender(ctx context.Context, obj *prisma.Administrator) (*prisma.Gender, error) {
	if obj.Deleted {
		return nil, nil
	}

	return &obj.Gender, nil
}
