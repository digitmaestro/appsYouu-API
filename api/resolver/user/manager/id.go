package manager

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Manager) ID(ctx context.Context, obj *prisma.Manager) (string, error) {
	if obj.Deleted {
		return "", nil
	}

	return obj.ID, nil
}
