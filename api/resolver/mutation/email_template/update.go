package email_template

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *EmailTemplateMutation) UpdateEmailTemplate(
	ctx context.Context,
	input gqlgen.UpdateEmailTemplateInput,
	language *string,
) (*gqlgen.UpdateEmailTemplatePayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if viewer.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError("user is not allowed to change email template")
	}

	emailTemplate, err := r.Prisma.UpdateEmailTemplate(prisma.EmailTemplateUpdateParams{
		Where: prisma.EmailTemplateWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.EmailTemplateUpdateInput{
			Title:   i18n.UpdateRequiredLocalizedString(ctx, input.Patch.Title),
			Content: i18n.UpdateRequiredLocalizedString(ctx, input.Patch.Content),
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpdateEmailTemplatePayload{
		Template: emailTemplate,
	}, nil
}
