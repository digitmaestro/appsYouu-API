package sessctx

import (
	"context"
)

const CompanyHeaderContextKey = "companyHeader"

// SetCompanyHeader returns a context that includes the user value.
func SetCompanyHeader(ctx context.Context, companyHeader string) context.Context {
	return context.WithValue(ctx, CompanyHeaderContextKey, companyHeader)
}

// CompanyHeader returns the company header from the context.
func CompanyHeader(ctx context.Context) string {
	return ctx.Value(CompanyHeaderContextKey).(string)
}
