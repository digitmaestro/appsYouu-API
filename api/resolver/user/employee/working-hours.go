package employee

import (
	"context"
	"time"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/hours"
	"github.com/steebchen/keskin-api/prisma"
)

var berlin, _ = time.LoadLocation("Europe/Berlin")

func (r *Employee) WorkingHours(ctx context.Context, obj *prisma.Employee) (*gqlgen.WorkingHours, error) {
	if obj.Deleted {
		return nil, nil
	}

	workingHours, err := r.Prisma.User(prisma.UserWhereUniqueInput{
		ID: &obj.ID,
	}).WorkingHours(nil).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.WorkingHours{
		Formatted: hours.FormatWorkingHoursPtr(hours.FormatWorkingHours(ctx, berlin, workingHours)),
		Raw:       hours.RawWorkingHoursPtr(hours.RawWorkingHours(ctx, berlin, workingHours)),
	}, nil
}
