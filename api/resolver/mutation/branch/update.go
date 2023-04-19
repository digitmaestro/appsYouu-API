package branch

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/file"
	"github.com/steebchen/keskin-api/lib/mailchimp"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *BranchMutation) UpdateBranch(
	ctx context.Context,
	input gqlgen.UpdateBranchInput,
	language *string,
) (*gqlgen.UpdateBranchPayload, error) {
	if err := permissions.CanAccessBranch(ctx, input.ID, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	imageID, err := file.MaybeUpload(input.Patch.Image, true)

	if err != nil {
		return nil, err
	}

	logoID, err := file.MaybeUpload(input.Patch.Logo, false)

	if err != nil {
		return nil, err
	}

	mailchimpApiKey := input.Patch.MailchimpAPIKey
	mailchimpListId := input.Patch.MailchimpListID

	branch, err := r.Prisma.UpdateBranch(prisma.BranchUpdateParams{
		Where: prisma.BranchWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.BranchUpdateInput{
			Name:               i18n.UpdateRequiredLocalizedString(ctx, input.Patch.Name),
			PhoneNumber:        input.Patch.PhoneNumber,
			Address:            input.Patch.Address,
			WelcomeMessage:     i18n.UpdateRequiredLocalizedString(ctx, input.Patch.WelcomeMessage),
			Image:              imageID,
			Logo:               logoID,
			AppTheme:           input.Patch.AppTheme,
			SmtpSendHost:       input.Patch.SMTPSendHost,
			SmtpSendPort:       input.Patch.SMTPSendPort,
			SmtpUsername:       input.Patch.SMTPUsername,
			SmtpPassword:       input.Patch.SMTPPassword,
			FromEmail:          input.Patch.FromEmail,
			WebsiteUrl:         input.Patch.WebsiteURL,
			NavigationLink:     input.Patch.NavigationLink,
			SharingRedirectUrl: input.Patch.SharingRedirectURL,
			MailchimpApiKey:    mailchimpApiKey,
			MailchimpListId:    mailchimpListId,
			Imprint:            input.Patch.Imprint,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if mailchimpApiKey != nil && mailchimpListId != nil {
		mailchimp.AssertMergeFields(*mailchimpApiKey, *mailchimpListId)
	}

	return &gqlgen.UpdateBranchPayload{
		Branch: branch,
	}, nil
}
