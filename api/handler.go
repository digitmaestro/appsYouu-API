package api

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"

	"github.com/steebchen/keskin-api/api/resolver/root"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/auth"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type Handler struct {
	Next http.Handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Next.ServeHTTP(w, r)
}

type HandlerFuncAdapter struct {
	NextFunc http.HandlerFunc
}

func (h *HandlerFuncAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.NextFunc(w, r)
}

func parseLanguage(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	resolverContext := graphql.GetResolverContext(ctx)
	language, ok := resolverContext.Args["language"]
	if ok && language != nil {
		ctx = sessctx.SetLanguage(ctx, language.(*string))
	}

	return next(ctx)
}

func handleLogin(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	loginFix := sessctx.LoginFix(ctx)

	if loginFix.Apply {
		ctx = sessctx.SetToken(ctx, loginFix.SessionToken)
		ctx = sessctx.SetUser(ctx, loginFix.User)
	}

	return next(ctx)
}

func prepareLoginFix(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	ctx = sessctx.InitLoginFix(ctx)

	return next(ctx)
}

func New(client *prisma.Client, resolver *root.Resolver) *Handler {
	schema := gqlgen.NewExecutableSchema(gqlgen.Config{Resolvers: resolver})

	return &Handler{
		Next: &auth.Handler{
			Prisma: client,
			Next: &HandlerFuncAdapter{
				NextFunc: handler.GraphQL(schema, handler.ResolverMiddleware(parseLanguage), handler.ResolverMiddleware(handleLogin), handler.RequestMiddleware(prepareLoginFix)),
			},
		},
	}
}
