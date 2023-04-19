package manager

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Manager) Gender(ctx context.Context, obj *prisma.Manager) (*prisma.Gender, error) {
	if obj.Deleted {
		return nil, nil
	}

	return &obj.Gender, nil
}
