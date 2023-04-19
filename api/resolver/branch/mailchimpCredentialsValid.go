package branch

import (
	"context"

	"github.com/steebchen/keskin-api/prisma"
)

func (r *Branch) MailchimpCredentialsValid(ctx context.Context, obj *prisma.Branch) (*bool, error) {
	valid := obj.MailchimpApiKey != nil && obj.MailchimpListId != nil

	return &valid, nil
}
