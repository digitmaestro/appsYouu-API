package email_template

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *EmailTemplateMutation) SendInviteEmail(
	ctx context.Context,
	input gqlgen.SendInviteEmailInput,
	language *string,
) (*gqlgen.SendEmailPayload, error) {
	viewer, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if viewer.Type != prisma.UserTypeManager && viewer.Type != prisma.UserTypeEmployee && viewer.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError("user is not allowed to send mails")
	}

	if err := permissions.CanAccessBranch(ctx, input.Branch, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	return SendEmailTemplate(
		ctx,
		r.Prisma,
		"invite",
		input.Branch,
		input.Email,
		input.Gender,
		"",
		input.Name,
		nil,
		nil,
		nil,
		nil,
	)
}
