package email_template

import (
	"context"

	"github.com/cbroglie/mustache"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/share"
	"github.com/steebchen/keskin-api/lib/strings"
	"github.com/steebchen/keskin-api/prisma"
)

func SendEmailTemplate(
	ctx context.Context,
	prismaClient *prisma.Client,
	templateName string,
	branchId string,
	toEmail string,
	gender prisma.Gender,
	lastName string,
	firstName string,
	appointmentDate *string,
	appointmentTime *string,
	passwordToken *string,
	activateToken *string,
) (*gqlgen.SendEmailPayload, error) {
	company, err := prismaClient.Branch(prisma.BranchWhereUniqueInput{
		ID: &branchId,
	}).Company().Exec(ctx)

	if err != nil {
		return nil, err
	}

	companyName, err := prismaClient.Branch(prisma.BranchWhereUniqueInput{
		ID: &branchId,
	}).Company().Name().Exec(ctx)

	if err != nil {
		return nil, err
	}

	templateContent, err := prismaClient.EmailTemplate(prisma.EmailTemplateWhereUniqueInput{
		Name: &templateName,
	}).Content().Exec(ctx)

	if err != nil {
		return nil, err
	}

	templateTitle, err := prismaClient.EmailTemplate(prisma.EmailTemplateWhereUniqueInput{
		Name: &templateName,
	}).Title().Exec(ctx)

	if err != nil {
		return nil, err
	}

	saluation := ""

	if gender == prisma.GenderMale {
		saluation = i18n.Language(ctx)["SALUTATION_MALE"]
	} else {
		saluation = i18n.Language(ctx)["SALUTATION_FEMALE"]
	}

	templateParameters := map[string]string{
		"salutation":        saluation,
		"lastName":          lastName,
		"firstName":         firstName,
		"companyName":       strings.DefaultWhenEmpty(i18n.GetLocalizedString(ctx, companyName), "appsYouu"),
		"customerSubdomain": share.ResolveCompanyUrl(ctx, prismaClient, company.ID),
	}

	if appointmentDate != nil {
		templateParameters["appointmentDate"] = *appointmentDate
	}

	if appointmentTime != nil {
		templateParameters["appointmentTime"] = *appointmentTime
	}

	if passwordToken != nil {
		templateParameters["passwordToken"] = *passwordToken
	}

	if activateToken != nil {
		templateParameters["activateToken"] = *activateToken
	}

	filledTemplate, err := mustache.Render(strings.DefaultWhenEmpty(i18n.GetLocalizedString(ctx, templateContent), ""), templateParameters)

	if err != nil {
		return nil, err
	}

	return SendEmail(
		ctx,
		prismaClient,
		branchId,
		toEmail,
		strings.DefaultWhenEmpty(i18n.GetLocalizedString(ctx, templateTitle), ""),
		filledTemplate,
	)
}
