package query

import (
	"context"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Query) CurrentCompany(ctx context.Context) (*prisma.Company, error) {
	ctxCompanyID := sessctx.CompanyWithFallback(ctx, r.Prisma, nil)

	company, err := r.Prisma.Company(prisma.CompanyWhereUniqueInput{
		ID: &ctxCompanyID,
	}).Exec(ctx)

	if err != nil && err != prisma.ErrNoResult {
		return company, err
	}

	return company, nil
}
