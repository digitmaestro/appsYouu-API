package mailchimp

import (
	"github.com/steebchen/keskin-api/prisma"
)

type MailchimpMutation struct {
	Prisma *prisma.Client
}

func New(client *prisma.Client) *MailchimpMutation {
	return &MailchimpMutation{
		Prisma: client,
	}
}
