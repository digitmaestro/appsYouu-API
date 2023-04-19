package sessctx

import (
	"context"
	"strings"
)

const LanguageContextKey = "language"

func SetLanguage(ctx context.Context, language *string) context.Context {
	if language != nil {
		languageKey := strings.ToUpper(*language)
		return context.WithValue(ctx, LanguageContextKey, &languageKey)
	} else {
		return ctx
	}
}

// Language returns the language from the context.
func Language(ctx context.Context) *string {
	value := ctx.Value(LanguageContextKey)

	if value != nil {
		return value.(*string)
	} else {
		return nil
	}
}
