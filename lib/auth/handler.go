package auth

import (
	"net/http"
	"strings"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type Handler struct {
	Prisma *prisma.Client
	Next   http.Handler
}

const CompanyHeader = "X-Company"
const LanguageHeader = "X-Language"

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Allowing all of the hosts
	origin := r.Header.Get("Host")
	w.Header().Set("Access-Control-Allow-Origin", origin)

	ctx := sessctx.SetWriter(r.Context(), w)

	ctx = sessctx.SetCompanyHeader(ctx, r.Header.Get(CompanyHeader))
	ctx = sessctx.SetHost(ctx, r.Host)

	language := r.Header.Get(LanguageHeader)
	ctx = sessctx.SetLanguage(ctx, &language)

	authorizationHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authorizationHeader, "Bearer ")

	if len(splitToken) <= 1 {
		h.Next.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	token := splitToken[1]

	user, err := h.Prisma.Session(prisma.SessionWhereUniqueInput{
		Token: &token,
	}).User().Exec(r.Context())

	if err == prisma.ErrNoResult || (user != nil && (user.Deleted || !user.Activated)) {
		// session removed or invalid
		h.Next.ServeHTTP(w, r.WithContext(ctx))
		return
	} else if err != nil {
		panic(err)
	}

	user, err = h.Prisma.UpdateUser(prisma.UserUpdateParams{
		Where: prisma.UserWhereUniqueInput{
			ID: &user.ID,
		},
		Data: prisma.UserUpdateInput{
			Language: &language,
		},
	}).Exec(ctx)

	if err != nil {
		panic(err)
	}

	company, err := h.Prisma.Session(prisma.SessionWhereUniqueInput{
		Token: &token,
	}).Company().Exec(r.Context())

	if err != nil && err != prisma.ErrNoResult {
		panic(err)
	}

	ctx = sessctx.SetToken(ctx, token)
	ctx = sessctx.SetUser(ctx, user)
	ctx = sessctx.SetCompany(ctx, company)

	h.Next.ServeHTTP(w, r.WithContext(ctx))
}
