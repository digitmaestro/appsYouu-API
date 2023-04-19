package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

const CookieKey = "session"

// SetCookie the cookie for the session
func SetCookie(ctx context.Context, session *prisma.Session) {
	w := sessctx.Writer(ctx)

	cookie := &http.Cookie{
		Value:    session.Token,
		Name:     CookieKey,
		HttpOnly: true,
		// TODO: Secure: env.Env == env.Production,
		Expires:  time.Now().AddDate(1, 0, 0),
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
}

// UnsetCookie the session cookie
func UnsetCookie(ctx context.Context) {
	w := sessctx.Writer(ctx)

	cookie := &http.Cookie{
		Name:   CookieKey,
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}
