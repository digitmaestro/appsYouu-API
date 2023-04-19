package share

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

var DefaultShareRedirectUrl string = ""
var BaseHostName string = ""

func ResolveCompanyUrl(ctx context.Context, prismaClient *prisma.Client, companyId string) string {
	var url *string = nil

	customUrls, err := prismaClient.CustomUrls(&prisma.CustomUrlsParams{
		Where: &prisma.CustomUrlWhereInput{
			Company: &prisma.CompanyWhereInput{
				ID: &companyId,
			},
		},
	}).Exec(ctx)

	if err == nil && len(customUrls) > 0 {
		customUrl := "https://" + customUrls[0].Value
		url = &customUrl
	}

	if url == nil {
		aliases, err := prismaClient.Aliases(&prisma.AliasesParams{
			Where: &prisma.AliasWhereInput{
				Company: &prisma.CompanyWhereInput{
					ID: &companyId,
				},
			},
		}).Exec(ctx)

		if err == nil && len(aliases) > 0 {
			aliasUrl := "https://" + aliases[0].Value + "." + BaseHostName
			url = &aliasUrl
		}
	}

	if url == nil {
		companyIdUrl := "https://" + companyId + "." + BaseHostName
		url = &companyIdUrl
	}

	return *url
}

func ResolveCompanyUrlFromBranchId(ctx context.Context, prismaClient *prisma.Client, branchId string) string {
	company, err := prismaClient.Branch(prisma.BranchWhereUniqueInput{
		ID: &branchId,
	}).Company().Exec(ctx)

	if err != nil {
		panic(err)
	}

	return ResolveCompanyUrl(ctx, prismaClient, company.ID)
}

func ResolveShareRedirectUrl(ctx context.Context, prismaClient *prisma.Client, branch *prisma.Branch) string {
	company, err := prismaClient.Branch(prisma.BranchWhereUniqueInput{
		ID: &branch.ID,
	}).Company().Exec(ctx)

	url := branch.SharingRedirectUrl

	if err == nil && company != nil {
		if url == nil {
			url = company.SharingRedirectUrl
		}

		if url == nil {
			companyUrl := ResolveCompanyUrl(ctx, prismaClient, company.ID)
			url = &companyUrl
		}
	}

	if url == nil {
		return DefaultShareRedirectUrl
	}

	return *url
}
