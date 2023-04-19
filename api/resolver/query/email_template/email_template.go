package email_template

import (
	"context"

	"github.com/google/wire"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type EmailTemplateQuery struct {
	Prisma *prisma.Client
}

func (r *EmailTemplateQuery) EmailTemplate(ctx context.Context, input gqlgen.EmailTemplateInput, language *string) (*gqlgen.EmailTemplateQueryPayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if viewer.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError("user is not allowed to access email templates")
	}

	where := prisma.EmailTemplateWhereUniqueInput{}

	if input.ID != nil {
		where.ID = input.ID
	} else {
		where.Name = input.Name
	}

	emailTemplate, err := r.Prisma.EmailTemplate(where).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.EmailTemplateQueryPayload{
		Template: emailTemplate,
	}, nil
}

func (r *EmailTemplateQuery) EmailTemplates(ctx context.Context, language *string) (*gqlgen.EmailTemplatesQueryPayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if viewer.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError("user is not allowed to access email templates")
	}

	emailTemplates, err := r.Prisma.EmailTemplates(&prisma.EmailTemplatesParams{}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	templates := []*prisma.EmailTemplate{}

	for _, emailTemplate := range emailTemplates {
		clone := emailTemplate
		templates = append(templates, &clone)
	}

	return &gqlgen.EmailTemplatesQueryPayload{
		Templates: templates,
	}, nil
}

func New(client *prisma.Client) *EmailTemplateQuery {
	return &EmailTemplateQuery{
		Prisma: client,
	}
}

var ProviderSet = wire.NewSet(
	New,
)
