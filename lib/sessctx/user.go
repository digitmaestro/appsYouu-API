package sessctx

import (
	"context"

	"github.com/vektah/gqlparser/v2/gqlerror"

	"github.com/steebchen/keskin-api/prisma"
)

const UserContextKey = "user"

var UserNotLoggedInError = &gqlerror.Error{
	Message: "user not logged in",
	Extensions: map[string]interface{}{
		"type": "Auth",
		"name": "NotLoggedIn",
	},
}

// SetUser returns a context that includes the user value.
func SetUser(ctx context.Context, user *prisma.User) context.Context {
	return context.WithValue(ctx, UserContextKey, user)
}

// User returns the user value from the context.
// If the user value is missing UserNotLoggedInError is returned.
func User(ctx context.Context) (*prisma.User, error) {
	user, ok := ctx.Value(UserContextKey).(*prisma.User)

	if !ok || user.Deleted || !user.Activated {
		return nil, UserNotLoggedInError
	}

	return user, nil
}
