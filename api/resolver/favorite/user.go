package favorite

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Favorite) User(ctx context.Context, obj *prisma.Favorite) (prisma.IUser, error) {
	user, err := r.Prisma.Favorite(prisma.FavoriteWhereUniqueInput{
		ID: &obj.ID,
	}).User().Exec(ctx)

	if err != nil {
		return nil, err
	}

	iuser := user.Convert()

	return iuser, err
}
