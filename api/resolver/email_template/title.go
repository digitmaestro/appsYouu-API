package email_template

import (
	"context"

	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *EmailTemplate) Title(ctx context.Context, obj *prisma.EmailTemplate) (*string, error) {
	title, err := r.Prisma.EmailTemplate(prisma.EmailTemplateWhereUniqueInput{
		ID: &obj.ID,
	}).Title().Exec(ctx)

	if err == prisma.ErrNoResult {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return i18n.GetLocalizedString(ctx, title), err
}
