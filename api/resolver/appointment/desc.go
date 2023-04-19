package appointment

import (
	"context"

	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

func (r *Appointment) Desc(ctx context.Context, obj *prisma.Appointment) (string, error) {
	desc, err := r.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &obj.ID,
	}).Desc().Exec(ctx)

	localizedValue := i18n.GetLocalizedString(ctx, desc)

	if localizedValue == nil {
		return "", err
	} else {
		return *localizedValue, err
	}
}
