package appointment

import (
	"context"
	"fmt"
	"image"
	"image/draw"

	"github.com/disintegration/imaging"
	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/gqlgen/gqlerrors"
	"github.com/steebchen/keskin-api/lib/file"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/prisma"
)

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func (r *Appointment) UpdateCombinedImage(ctx context.Context, appointment *prisma.Appointment) error {
	if appointment.BeforeImage != nil && appointment.AfterImage != nil {
		beforeImagePath := file.Path(*appointment.BeforeImage)
		afterImagePath := file.Path(*appointment.AfterImage)

		beforeImage, err := picture.LoadImage(beforeImagePath)

		if err != nil {
			return err
		}

		afterImage, err := picture.LoadImage(afterImagePath)

		if err != nil {
			return err
		}

		beforeBounds := beforeImage.Bounds()
		afterBounds := afterImage.Bounds()

		bounds := image.Rect(beforeBounds.Min.X, beforeBounds.Min.Y, beforeBounds.Max.X+afterBounds.Max.X-afterBounds.Min.X, max(beforeBounds.Max.Y, afterBounds.Max.Y))
		combinedImage := image.NewRGBA(bounds)

		afterBounds.Min.X += beforeBounds.Max.X
		afterBounds.Max.X += beforeBounds.Max.X

		draw.Draw(combinedImage, beforeBounds, beforeImage, image.ZP, draw.Src)
		draw.Draw(combinedImage, afterBounds, afterImage, image.ZP, draw.Src)

		combinedImageId := fmt.Sprintf("%s_%s.jpg", picture.IDFromFileName(*appointment.BeforeImage), picture.IDFromFileName(*appointment.AfterImage))
		combinedImagePath := file.Path(combinedImageId)

		scaledImage := file.ScaleToFit(combinedImage)
		imaging.Save(scaledImage, combinedImagePath, imaging.JPEGQuality(80))
	}

	return nil
}

func (r *Appointment) UpdateBeforeImage(ctx context.Context, input gqlgen.UpdateAppointmentImageInput, language *string) (*gqlgen.UpdateAppointmentImagePayload, error) {
	authorized, err := r.MayUserModifyAppointment(ctx, input.Appointment)

	if err != nil {
		return nil, err
	}

	if !authorized {
		return nil, gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	imageID, err := file.MaybeUpload(&input.Image, true)

	if err != nil {
		return nil, err
	}

	appointment, err := r.Prisma.UpdateAppointment(prisma.AppointmentUpdateParams{
		Where: prisma.AppointmentWhereUniqueInput{
			ID: &input.Appointment,
		},
		Data: prisma.AppointmentUpdateInput{
			BeforeImage: imageID,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	err = r.UpdateCombinedImage(ctx, appointment)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpdateAppointmentImagePayload{
		Appointment: appointment,
	}, nil
}

func (r *Appointment) UpdateAfterImage(ctx context.Context, input gqlgen.UpdateAppointmentImageInput, language *string) (*gqlgen.UpdateAppointmentImagePayload, error) {
	authorized, err := r.MayUserModifyAppointment(ctx, input.Appointment)

	if err != nil {
		return nil, err
	}

	if !authorized {
		return nil, gqlerrors.NewPermissionError(gqlerrors.InvalidUserType)
	}

	imageID, err := file.MaybeUpload(&input.Image, true)

	if err != nil {
		return nil, err
	}

	appointment, err := r.Prisma.UpdateAppointment(prisma.AppointmentUpdateParams{
		Where: prisma.AppointmentWhereUniqueInput{
			ID: &input.Appointment,
		},
		Data: prisma.AppointmentUpdateInput{
			AfterImage: imageID,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	err = r.UpdateCombinedImage(ctx, appointment)

	if err != nil {
		return nil, err
	}

	return &gqlgen.UpdateAppointmentImagePayload{
		Appointment: appointment,
	}, nil
}
