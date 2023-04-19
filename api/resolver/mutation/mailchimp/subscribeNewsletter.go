package mailchimp

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"strings"

	"github.com/hanzoai/gochimp3"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/mailchimp"
	"github.com/steebchen/keskin-api/prisma"
)

type Address struct {
	Address1   string `json:"addr1"`
	City       string `json:"city"`
	PostalCode string `json:"zip"`
	State      string `json:"state"`
	Country    string `json:"country"`
}

func stringOrDefault(s *string, d string) string {
	if s != nil {
		return *s
	} else {
		return d
	}
}

func (r *MailchimpMutation) SubscribeNewsletter(
	ctx context.Context,
	email string,
	branchId string,
) (*gqlgen.SubscribeNewsletterPayload, error) {
	branch, err := r.Prisma.Branch(prisma.BranchWhereUniqueInput{
		ID: &branchId,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	var user *prisma.User = nil
	deleted := false

	users, err := r.Prisma.Users(&prisma.UsersParams{
		First: prisma.Int32(1),
		Where: &prisma.UserWhereInput{
			Branch: &prisma.BranchWhereInput{
				ID: &branchId,
			},
			Email:   &email,
			Deleted: &deleted,
		},
	}).Exec(ctx)

	if err != nil && err != prisma.ErrNoResult {
		return nil, err
	}

	if len(users) > 0 {
		user = &users[0]
	}

	if user == nil {
		users, err = r.Prisma.Users(&prisma.UsersParams{
			First: prisma.Int32(1),
			Where: &prisma.UserWhereInput{
				Company: &prisma.CompanyWhereInput{
					BranchesSome: &prisma.BranchWhereInput{
						ID: &branchId,
					},
				},
				Email:   &email,
				Deleted: &deleted,
			},
		}).Exec(ctx)

		if err != nil && err != prisma.ErrNoResult {
			return nil, err
		}

		if len(users) > 0 {
			user = &users[0]
		}
	}

	if user == nil {
		users, err = r.Prisma.Users(&prisma.UsersParams{
			First: prisma.Int32(1),
			Where: &prisma.UserWhereInput{
				Email:   &email,
				Deleted: &deleted,
			},
		}).Exec(ctx)

		if err != nil && err != prisma.ErrNoResult {
			return nil, err
		}

		if len(users) > 0 {
			user = &users[0]
		}
	}

	mailchimpApiKey := branch.MailchimpApiKey
	mailchimpListId := branch.MailchimpListId

	if mailchimpApiKey == nil || mailchimpListId == nil {
		return nil, errors.New("Mailchimp configuration incomplete")
	}

	client := gochimp3.New(*mailchimpApiKey)
	list, err := client.GetList(*mailchimpListId, nil)

	if err != nil {
		return nil, err
	}

	mergeFields := make(map[string]interface{})
	if user != nil {
		birthday := ""
		if user.Birthday != nil {
			birthday = (*user.Birthday)[5:7] + "/" + (*user.Birthday)[8:10]
		}

		mergeFields[mailchimp.UserIdFieldId] = user.ID
		mergeFields["FNAME"] = user.FirstName
		mergeFields["LNAME"] = user.LastName
		mergeFields["ADDRESS"] = Address{
			Address1:   stringOrDefault(user.Street, ""),
			PostalCode: stringOrDefault(user.ZipCode, ""),
			City:       stringOrDefault(user.City, ""),
			State:      "",
			Country:    "Germany",
		}
		mergeFields["BIRTHDAY"] = birthday
		mergeFields["PHONE"] = stringOrDefault(user.PhoneNumber, "")
	}

	subscriberHash := fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(email))))

	member, err := list.AddOrUpdateMember(subscriberHash, &gochimp3.MemberRequest{
		EmailAddress: email,
		Status:       "pending",
		MergeFields:  mergeFields,
	})

	if err != nil {
		return nil, err
	}

	return &gqlgen.SubscribeNewsletterPayload{
		ID:            &member.ID,
		UniqueEmailID: &member.UniqueEmailID,
	}, nil
}
