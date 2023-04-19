package administrator

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/user/fullname"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/prisma"
)

type Administrator struct {
	Prisma *prisma.Client
}

func (r *Administrator) FullName(ctx context.Context, obj *prisma.Administrator) (*string, error) {
	if obj.Deleted {
		deleted := i18n.Language(ctx)["DELETED_USER"]
		return &deleted, nil
	}

	return fullname.Convert(obj.FirstName, obj.LastName), nil
}

func (r *Administrator) Image(ctx context.Context, obj *prisma.Administrator) (*gqlgen.Image, error) {
	if obj.Deleted {
		return nil, nil
	}

	return picture.FromID(obj.Image), nil
}

func New(client *prisma.Client) *Administrator {
	return &Administrator{
		Prisma: client,
	}
}
