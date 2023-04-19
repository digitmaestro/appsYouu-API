package sessctx

import (
	"context"
	"strings"

	"github.com/steebchen/keskin-api/prisma"
)

const CompanyContextKey = "company"

// SetCompany returns a context that includes the user value.
func SetCompany(ctx context.Context, company *prisma.Company) context.Context {
	companyID := ""

	if company != nil {
		companyID = company.ID
	}

	return context.WithValue(ctx, CompanyContextKey, companyID)
}

// Company returns the company ID from the context.
func Company(ctx context.Context) string {
	value := ctx.Value(CompanyContextKey)
	if value != nil {
		return value.(string)
	} else {
		return ""
	}
}

func CompanyWithFallback(ctx context.Context, client *prisma.Client, input *string) string {
	splitHost := strings.Split(Host(ctx), ":")
	domainParts := strings.Split(splitHost[0], ".")

	companyID := ""

	if input != nil {
		companyID = *input
	}

	if companyID == "" {
		companyID = Company(ctx)
	}

	if companyID == "" {
		company, err := client.Company(prisma.CompanyWhereUniqueInput{
			ID: &domainParts[0],
		}).Exec(ctx)

		if err != nil && err != prisma.ErrNoResult {
			panic(err)
		}

		if company != nil {
			companyID = company.ID
		}
	}

	if companyID == "" {
		company, err := client.Alias(prisma.AliasWhereUniqueInput{
			Value: &domainParts[0],
		}).Company().Exec(ctx)

		if err != nil && err != prisma.ErrNoResult {
			panic(err)
		}

		if company != nil {
			companyID = company.ID
		}
	}

	if companyID == "" {
		company, err := client.CustomUrl(prisma.CustomUrlWhereUniqueInput{
			Value: &splitHost[0],
		}).Company().Exec(ctx)

		if err != nil && err != prisma.ErrNoResult {
			panic(err)
		}

		if company != nil {
			companyID = company.ID
		}
	}

	companyIDFromHeader := CompanyHeader(ctx)
	if companyIDFromHeader != "" {
		companyID = companyIDFromHeader
	}

	if companyID == "" && domainParts[0] == "localhost" {
		// set a default company when no other company was set in development mode
		companies, err := client.Companies(nil).Exec(ctx)

		if err != nil && err != prisma.ErrNoResult {
			panic(err)
		}

		if len(companies) > 0 {
			companyID = companies[0].ID
		}
	}

	return companyID
}
