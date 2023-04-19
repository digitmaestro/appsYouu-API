package company

import (
	"context"

	"github.com/steebchen/keskin-api/api/permissions"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/file"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *CompanyMutation) UpdateCompany(
	ctx context.Context,
	input gqlgen.UpdateCompanyInput,
	language *string,
) (*gqlgen.UpdateCompanyPayload, error) {
	err := permissions.CanAccessCompany(ctx, input.ID, r.Prisma, []prisma.UserType{prisma.UserTypeManager})

	if err != nil {
		return nil, err
	}

	var createCustomUrls *prisma.CustomUrlUpdateManyWithoutCompanyInput = nil
	var createAliases *prisma.AliasUpdateManyWithoutCompanyInput = nil

	if input.Patch.CustomUrls != nil {
		customUrls, err := r.Prisma.CustomUrls(&prisma.CustomUrlsParams{
			Where: &prisma.CustomUrlWhereInput{
				ValueIn: input.Patch.CustomUrls,
				Company: &prisma.CompanyWhereInput{
					IDNot: &input.ID,
				},
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

		_, err = r.Prisma.DeleteManyCustomUrls(&prisma.CustomUrlWhereInput{
			Company: &prisma.CompanyWhereInput{
				ID: &input.ID,
			},
		}).Exec(ctx)

		if err != nil {
			return nil, err
		}

		createCustomUrls = &prisma.CustomUrlUpdateManyWithoutCompanyInput{}

		for _, customUrl := range input.Patch.CustomUrls {
			createCustomUrls.Create = append(createCustomUrls.Create, prisma.CustomUrlCreateWithoutCompanyInput{
				Value: customUrl,
			})
		}
	}

	if input.Patch.Aliases != nil {
		aliases, err := r.Prisma.Aliases(&prisma.AliasesParams{
			Where: &prisma.AliasWhereInput{
				ValueIn: input.Patch.Aliases,
				Company: &prisma.CompanyWhereInput{
					IDNot: &input.ID,
				},
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

		_, err = r.Prisma.DeleteManyAliases(&prisma.AliasWhereInput{
			Company: &prisma.CompanyWhereInput{
				ID: &input.ID,
			},
		}).Exec(ctx)

		if err != nil {
			return nil, err
		}

		createAliases = &prisma.AliasUpdateManyWithoutCompanyInput{}

		for _, alias := range input.Patch.Aliases {
			createAliases.Create = append(createAliases.Create, prisma.AliasCreateWithoutCompanyInput{
				Value: alias,
			})
		}
	}

	logoID, err := file.MaybeUpload(input.Patch.Logo, false)

	if err != nil {
		return nil, err
	}

	pwaIconID, err := file.MaybeUpload(input.Patch.PwaIcon, false)

	if err != nil {
		return nil, err
	}

	picture.CreatePwaIconSizes(pwaIconID)

	company, err := r.Prisma.UpdateCompany(prisma.CompanyUpdateParams{
		Where: prisma.CompanyWhereUniqueInput{
			ID: &input.ID,
		},
		Data: prisma.CompanyUpdateInput{
			Name:               i18n.UpdateRequiredLocalizedString(ctx, input.Patch.Name),
			SharingRedirectUrl: input.Patch.SharingRedirectURL,
			CustomUrls:         createCustomUrls,
			Aliases:            createAliases,
			Logo:               logoID,
			AppTheme:           input.Patch.AppTheme,
			PwaShortName:       i18n.UpdateRequiredLocalizedString(ctx, input.Patch.PwaShortName),
			PwaIcon:            pwaIconID,
			PwaThemeColor:      input.Patch.PwaThemeColor,
			PwaBackgroundColor: input.Patch.PwaBackgroundColor,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpdateCompanyPayload{
		Company: company,
	}, nil
}
