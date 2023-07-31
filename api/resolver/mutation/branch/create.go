package branch

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/file"
	"github.com/steebchen/keskin-api/lib/mailchimp"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *BranchMutation) CreateBranch(
	ctx context.Context,
	input gqlgen.CreateBranchInput,
	language *string,
) (*gqlgen.CreateBranchPayload, error) {
	if err := permissions.CanAccessCompany(ctx, input.Company, r.Prisma, allowedTypes); err != nil {
		return nil, err
	}

	var imgIds []string = []string{}

	for _, img := range input.Data.Images {
		id, err := file.MaybeUpload(img, true)
		if err != nil {
			return nil, gqlerrors.NewValidationError("error uploading pictures", "ErrorUploadingPictures")
		}

		imgIds = append(imgIds, *id)
	}

	// imageID, err := file.MaybeUpload(input.Data.Image, true)

	// if err != nil {
	// 	return nil, err
	// }

	logoID, err := file.MaybeUpload(input.Data.Logo, false)

	if err != nil {
		return nil, err
	}

	mailchimpApiKey := input.Data.MailchimpAPIKey
	mailchimpListId := input.Data.MailchimpListID

	branch, err := r.Prisma.CreateBranch(prisma.BranchCreateInput{
		Name:           *i18n.CreateLocalizedString(ctx, &input.Data.Name),
		PhoneNumber:    input.Data.PhoneNumber,
		Address:        input.Data.Address,
		WelcomeMessage: *i18n.CreateLocalizedString(ctx, input.Data.WelcomeMessage),
		Images: &prisma.BranchCreateimagesInput{
			Set: imgIds,
		},
		Logo:               logoID,
		AppTheme:           input.Data.AppTheme,
		FacebookLink:       input.Data.FacebookLink,
		TiktokLink:         input.Data.TiktokLink,
		InstagramLink:      input.Data.InstagramLink,
		SmtpSendHost:       input.Data.SMTPSendHost,
		SmtpSendPort:       input.Data.SMTPSendPort,
		SmtpUsername:       input.Data.SMTPUsername,
		SmtpPassword:       input.Data.SMTPPassword,
		FromEmail:          input.Data.FromEmail,
		WebsiteUrl:         input.Data.WebsiteURL,
		NavigationLink:     input.Data.NavigationLink,
		SharingRedirectUrl: input.Data.SharingRedirectURL,
		MailchimpApiKey:    mailchimpApiKey,
		MailchimpListId:    mailchimpListId,
		Imprint:            input.Data.Imprint,

		Company: prisma.CompanyCreateOneWithoutBranchesInput{
			Connect: &prisma.CompanyWhereUniqueInput{
				ID: &input.Company,
			},
		},

		OpeningHours: &prisma.BranchOpeningHourCreateManyWithoutBranchInput{
			Create: []prisma.BranchOpeningHourCreateWithoutBranchInput{},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if mailchimpApiKey != nil && mailchimpListId != nil {
		mailchimp.AssertMergeFields(*mailchimpApiKey, *mailchimpListId)
	}

	return &gqlgen.CreateBranchPayload{
		Branch: branch,
	}, nil
}
