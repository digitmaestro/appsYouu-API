package auth

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"

	"github.com/steebchen/keskin-api/api/resolver/mutation/email_template"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

type Auth struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *Auth {
	return &Auth{
		Prisma: client,
	}
}

func RequestPasswordReset(prismaClient *prisma.Client, ctx context.Context, email string, company *string) error {
	deleted := false

	where := &prisma.UserWhereInput{
		Email:   &email,
		Deleted: &deleted,
	}

	companyID := sessctx.CompanyWithFallback(ctx, prismaClient, company)
	userTypeAdministrator := prisma.UserTypeAdministrator

	if companyID != "" {
		_, err := prismaClient.Company(prisma.CompanyWhereUniqueInput{
			ID: &companyID,
		}).Exec(ctx)

		if err == nil {
			userTypeCustomer := prisma.UserTypeCustomer
			userTypeManager := prisma.UserTypeManager
			userTypeEmployee := prisma.UserTypeEmployee

			where.Or = []prisma.UserWhereInput{{
				Type: &userTypeEmployee,
				Branch: &prisma.BranchWhereInput{
					Company: &prisma.CompanyWhereInput{
						ID: &companyID,
					},
				},
			}, {
				TypeIn: []prisma.UserType{userTypeCustomer, userTypeManager},
				Company: &prisma.CompanyWhereInput{
					ID: &companyID,
				},
			}, {
				Type: &userTypeAdministrator,
			}}
		} else {
			where.Type = &userTypeAdministrator
		}
	} else {
		where.Type = &userTypeAdministrator
	}

	users, err := prismaClient.Users(&prisma.UsersParams{
		Where: where,
	}).Exec(ctx)

	if len(users) != 0 && err == nil {
		user := users[0]

		prismaClient.DeleteManyPasswordTokens(&prisma.PasswordTokenWhereInput{
			User: &prisma.UserWhereInput{
				ID: &user.ID,
			},
		}).Exec(ctx)

		passwordToken, err := prismaClient.CreatePasswordToken(prisma.PasswordTokenCreateInput{
			Token: uuid.New().String(),
			User: prisma.UserCreateOneWithoutPasswordTokenInput{
				Connect: &prisma.UserWhereUniqueInput{
					ID: &user.ID,
				},
			},
		}).Exec(ctx)

		if passwordToken != nil && err == nil {
			var branchesParams *prisma.BranchesParams = nil

			if companyID != "" {
				branchesParams = &prisma.BranchesParams{
					Where: &prisma.BranchWhereInput{
						Company: &prisma.CompanyWhereInput{
							ID: &companyID,
						},
					},
				}
			}

			branches, err := prismaClient.Branches(branchesParams).Exec(ctx)

			if err == nil && len(branches) > 0 {
				branch := branches[0]

				for _, b := range branches {
					if b.FromEmail != nil && b.SmtpUsername != nil && b.SmtpPassword != nil && b.SmtpSendHost != nil && b.SmtpSendPort != nil {
						branch = b
					}
				}

				_, err := email_template.SendEmailTemplate(context.Background(), prismaClient, "passwordReset", branch.ID, user.Email, user.Gender, user.LastName, user.FirstName, nil, nil, &passwordToken.Token, nil)
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}

func RequestActivationLink(
	prismaClient *prisma.Client,
	ctx context.Context,
	email string,
	company *string,
) error {
	deleted := false
	activated := false

	where := &prisma.UserWhereInput{
		Email:     &email,
		Deleted:   &deleted,
		Activated: &activated,
	}

	companyID := sessctx.CompanyWithFallback(ctx, prismaClient, company)
	userTypeAdministrator := prisma.UserTypeAdministrator

	if companyID != "" {
		_, err := prismaClient.Company(prisma.CompanyWhereUniqueInput{
			ID: &companyID,
		}).Exec(ctx)

		if err == nil {
			userTypeCustomer := prisma.UserTypeCustomer
			userTypeManager := prisma.UserTypeManager
			userTypeEmployee := prisma.UserTypeEmployee

			where.Or = []prisma.UserWhereInput{{
				Type: &userTypeEmployee,
				Branch: &prisma.BranchWhereInput{
					Company: &prisma.CompanyWhereInput{
						ID: &companyID,
					},
				},
			}, {
				TypeIn: []prisma.UserType{userTypeCustomer, userTypeManager},
				Company: &prisma.CompanyWhereInput{
					ID: &companyID,
				},
			}, {
				Type: &userTypeAdministrator,
			}}
		} else {
			where.Type = &userTypeAdministrator
		}
	} else {
		where.Type = &userTypeAdministrator
	}

	users, err := prismaClient.Users(&prisma.UsersParams{
		Where: where,
	}).Exec(ctx)

	if len(users) != 0 && err == nil {
		user := users[0]
		activateToken := ""

		if user.ActivateToken == nil {
			activateToken, err := GenerateActivateToken(prismaClient, ctx)

			if err != nil {
				return err
			}

			updatedUser, err := prismaClient.UpdateUser(prisma.UserUpdateParams{ //nolint
				Where: prisma.UserWhereUniqueInput{
					ID: &user.ID,
				},
				Data: prisma.UserUpdateInput{
					ActivateToken: &activateToken,
				},
			}).Exec(ctx)

			if err != nil {
				log.Printf("error updating user %v", err)

				return errors.New("error requesting link")
			}

			if updatedUser != nil {
				user = *updatedUser
			}
		} else {
			activateToken = *user.ActivateToken
		}

		if user.ActivateToken != nil && err == nil {
			var branchesParams *prisma.BranchesParams = nil

			if companyID != "" {
				branchesParams = &prisma.BranchesParams{
					Where: &prisma.BranchWhereInput{
						Company: &prisma.CompanyWhereInput{
							ID: &companyID,
						},
					},
				}
			}

			branches, err := prismaClient.Branches(branchesParams).Exec(ctx)

			if err == nil && len(branches) > 0 {
				branch := branches[0]

				for _, b := range branches {
					if b.FromEmail != nil && b.SmtpUsername != nil && b.SmtpPassword != nil && b.SmtpSendHost != nil && b.SmtpSendPort != nil {
						branch = b
					}
				}

				_, err := email_template.SendEmailTemplate(context.Background(), prismaClient, "activationLink", branch.ID, user.Email, user.Gender, user.LastName, user.FirstName, nil, nil, nil, &activateToken)
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}

func GenerateActivateToken(prismaClient *prisma.Client, ctx context.Context) (string, error) {
	activateToken := uuid.New().String()

	users, err := prismaClient.Users(&prisma.UsersParams{
		Where: &prisma.UserWhereInput{
			ActivateToken: &activateToken,
		},
	}).Exec(ctx)

	if err != nil {
		return "", err
	}

	for len(users) != 0 {
		activateToken := uuid.New().String()

		users, err = prismaClient.Users(&prisma.UsersParams{
			Where: &prisma.UserWhereInput{
				ActivateToken: &activateToken,
			},
		}).Exec(ctx)

		if err != nil {
			return "", err
		}
	}

	return activateToken, nil
}
