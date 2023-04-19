package branch

import (
	"context"

	"github.com/steebchen/keskin-api/api/resolver/phonenumber"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Branch) PhoneNumber(ctx context.Context, obj *prisma.Branch) (*gqlgen.PhoneNumber, error) {
	return phonenumber.Convert(obj.PhoneNumber), nil
}
