package company

import (
	"context"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/file"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *CompanyMutation) CreateCompany(
	ctx context.Context,
	input gqlgen.CreateCompanyInput,
	language *string,
) (*gqlgen.CreateCompanyPayload, error) {
	user, err := sessctx.User(ctx)

	if err != nil {
		return nil, err
	}

	if user.Type != prisma.UserTypeAdministrator {
		return nil, gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	customUrls, err := r.Prisma.CustomUrls(&prisma.CustomUrlsParams{
		Where: &prisma.CustomUrlWhereInput{
			ValueIn: input.Data.CustomUrls,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(customUrls) > 0 {
		usedUrls := ""
		for _, url := range customUrls {
			if usedUrls != "" {
				usedUrls += ", "
			}
			usedUrls += url.Value
		}
		return nil, gqlerrors.NewValidationError("URLs already in use: "+usedUrls, "DuplicateEntry")
	}

	aliases, err := r.Prisma.Aliases(&prisma.AliasesParams{
		Where: &prisma.AliasWhereInput{
			ValueIn: input.Data.Aliases,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(aliases) > 0 {
		usedAliases := ""
		for _, alias := range aliases {
			if usedAliases != "" {
				usedAliases += ", "
			}
			usedAliases += alias.Value
		}
		return nil, gqlerrors.NewValidationError("Aliases already in use: "+usedAliases, "DuplicateEntry")
	}

	logoID, err := file.MaybeUpload(input.Data.Logo, false)

	if err != nil {
		return nil, err
	}

	createCustomUrls := prisma.CustomUrlCreateManyWithoutCompanyInput{}

	for _, customUrl := range input.Data.CustomUrls {
		createCustomUrls.Create = append(createCustomUrls.Create, prisma.CustomUrlCreateWithoutCompanyInput{
			Value: customUrl,
		})
	}

	createAliases := prisma.AliasCreateManyWithoutCompanyInput{}

	for _, alias := range input.Data.Aliases {
		createAliases.Create = append(createAliases.Create, prisma.AliasCreateWithoutCompanyInput{
			Value: alias,
		})
	}

	pwaIconID, err := file.MaybeUpload(input.Data.PwaIcon, false)

	if err != nil {
		return nil, err
	}

	picture.CreatePwaIconSizes(pwaIconID)

	company, err := r.Prisma.CreateCompany(prisma.CompanyCreateInput{
		Name:               *i18n.CreateLocalizedString(ctx, input.Data.Name),
		SharingRedirectUrl: input.Data.SharingRedirectURL,
		CustomUrls:         &createCustomUrls,
		Aliases:            &createAliases,
		Logo:               logoID,
		AppTheme:           input.Data.AppTheme,
		PwaShortName:       *i18n.CreateLocalizedString(ctx, &input.Data.PwaShortName),
		PwaIcon:            pwaIconID,
		PwaThemeColor:      input.Data.PwaThemeColor,
		PwaBackgroundColor: input.Data.PwaBackgroundColor,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.CreateCompanyPayload{
		Company: company,
	}, nil
}
