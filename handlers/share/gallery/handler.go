package gallery

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cbroglie/mustache"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/lib/share"
	"github.com/steebchen/keskin-api/prisma"
)

type Handler struct {
	Prisma *prisma.Client
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<meta property="og:title" content="{{reviewTitle}}" />
		<meta property="og:description" content="{{reviewText}}" />
		<meta property="og:image" content="{{companyUrl}}{{branchImageUrl}}" />
		<meta property="og:site_name" content="{{siteName}}">

		<meta property="twitter:card" content="summary_large_image" />
		<meta property="twitter:title" content="{{reviewTitle}}" />
		<meta property="twitter:description" content="{{reviewText}}" />
		<meta property="twitter:image" content="{{companyUrl}}{{branchImageUrl}}" />

		<title>{{reviewTitle}}</title>

		<noscript>
			<meta http-equiv="refresh" content="0; url={{linkUrl}}" />
		</noscript>
	</head>
	<body>
		<script>
			setTimeout(function() {
				location.href = '{{linkUrl}}';
			}, 1000);
		</script>
	</body>
</html>
`

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	urlParts := strings.Split(r.URL.Path, "/")
	appointmentId := urlParts[len(urlParts)-1]

	appointments, err := h.Prisma.Appointments(&prisma.AppointmentsParams{
		Where: &prisma.AppointmentWhereInput{
			ID: &appointmentId,
		},
	}).Exec(ctx)

	if err != nil || len(appointments) == 0 {
		w.WriteHeader(404)
		return
	}

	appointment := appointments[0]

	customer, err := h.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &appointment.ID,
	}).Customer().Exec(ctx)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	if customer.Deleted || !customer.AllowReviewSharing {
		w.WriteHeader(404)
		return
	}

	branch, err := h.Prisma.Appointment(prisma.AppointmentWhereUniqueInput{
		ID: &appointment.ID,
	}).Branch().Exec(ctx)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	var image *gqlgen.Image = nil

	if appointment.AfterImage != nil && appointment.BeforeImage != nil {
		combinedImageId := fmt.Sprintf("%s_%s.jpg", picture.IDFromFileName(*appointment.BeforeImage), picture.IDFromFileName(*appointment.AfterImage))

		image = picture.FromID(&combinedImageId)
	}

	if image == nil {
		image = picture.FromID(appointment.AfterImage)
	}

	if image == nil {
		image = picture.FromID(appointment.BeforeImage)
	}

	templateParameters := map[string]interface{}{
		"reviewTitle":     "",
		"reviewText":      "",
		"linkUrl":         "",
		"branchImageUrls": []string{},
		"companyUrl":      share.ResolveCompanyUrlFromBranchId(ctx, h.Prisma, branch.ID),
		"siteName":        i18n.Language(ctx)["SITE_NAME"],
	}

	if branch != nil {
		branchName, _ := h.Prisma.Branch(prisma.BranchWhereUniqueInput{
			ID: &branch.ID,
		}).Name().Exec(ctx)

		localizedBranchName := i18n.GetLocalizedString(ctx, branchName)

		if localizedBranchName != nil {
			templateParameters["reviewTitle"] = *localizedBranchName
		}

		if image == nil {
			image = picture.FromID(&branch.Images[0])
		}

		if image != nil {
			templateParameters["branchImageUrl"] = image.URL
		}

		if branch.Address != nil {
			templateParameters["reviewText"] = *branch.Address
		}
		templateParameters["linkUrl"] = share.ResolveShareRedirectUrl(ctx, h.Prisma, branch)
	}

	filledTemplate, err := mustache.Render(htmlTemplate, templateParameters)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Write([]byte(filledTemplate))
}

func New(client *prisma.Client) *Handler {
	return &Handler{
		Prisma: client,
	}
}
