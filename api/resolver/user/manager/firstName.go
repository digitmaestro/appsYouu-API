package manager

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Manager) FirstName(ctx context.Context, obj *prisma.Manager) (*string, error) {
	if obj.Deleted {
		return nil, nil
	}

	return &obj.FirstName, nil
}
