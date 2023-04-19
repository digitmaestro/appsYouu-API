package review

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
	reviewId := urlParts[len(urlParts)-1]

	reviews, err := h.Prisma.Reviews(&prisma.ReviewsParams{
		Where: &prisma.ReviewWhereInput{
			ID: &reviewId,
		},
	}).Exec(ctx)

	if err != nil || len(reviews) == 0 {
		w.WriteHeader(404)
		return
	}

	review := reviews[0]

	customer, err := h.Prisma.Review(prisma.ReviewWhereUniqueInput{
		ID: &review.ID,
	}).Customer().Exec(ctx)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	if customer.Deleted || !customer.AllowReviewSharing {
		w.WriteHeader(404)
		return
	}

	var branch *prisma.Branch = nil
	var image *gqlgen.Image = nil
	reviewTitle := ""

	if review.Type == prisma.ReviewTypeProduct {
		branch, err = h.Prisma.Review(prisma.ReviewWhereUniqueInput{
			ID: &review.ID,
		}).Product().Branch().Exec(ctx)

		var product *prisma.Product
		var localizedProductName *string

		product, err = h.Prisma.Review(prisma.ReviewWhereUniqueInput{
			ID: &review.ID,
		}).Product().Exec(ctx)

		if product != nil && !product.Deleted {
			image = picture.FromID(product.Image)

			var productName *prisma.LocalizedString

			productName, err = h.Prisma.Review(prisma.ReviewWhereUniqueInput{
				ID: &review.ID,
			}).Product().Name().Exec(ctx)

			localizedProductName = i18n.GetLocalizedString(ctx, productName)
		} else {
			deleted := "Gelöschtes Produkt"
			localizedProductName = &deleted
		}

		if localizedProductName != nil {
			reviewTitle += *localizedProductName
		}
	} else if review.Type == prisma.ReviewTypeService {
		branch, err = h.Prisma.Review(prisma.ReviewWhereUniqueInput{
			ID: &review.ID,
		}).Service().Branch().Exec(ctx)

		var service *prisma.Service
		var localizedServiceName *string

		service, err = h.Prisma.Review(prisma.ReviewWhereUniqueInput{
			ID: &review.ID,
		}).Service().Exec(ctx)

		if service != nil && !service.Deleted {
			image = picture.FromID(service.Image)

			var serviceName *prisma.LocalizedString

			serviceName, err = h.Prisma.Review(prisma.ReviewWhereUniqueInput{
				ID: &review.ID,
			}).Service().Name().Exec(ctx)

			localizedServiceName = i18n.GetLocalizedString(ctx, serviceName)
		} else {
			deleted := "Gelöschte Dienstleistung"
			localizedServiceName = &deleted
		}

		if localizedServiceName != nil {
			reviewTitle += *localizedServiceName
		}
	} else if review.Type == prisma.ReviewTypeAppointment {
		branch, err = h.Prisma.Review(prisma.ReviewWhereUniqueInput{
			ID: &review.ID,
		}).Appointment().Branch().Exec(ctx)

		var appointment *prisma.Appointment

		appointment, err = h.Prisma.Review(prisma.ReviewWhereUniqueInput{
			ID: &review.ID,
		}).Appointment().Exec(ctx)

		if appointment != nil && appointment.AfterImage != nil && appointment.BeforeImage != nil {
			combinedImageId := fmt.Sprintf("%s_%s.jpg", picture.IDFromFileName(*appointment.BeforeImage), picture.IDFromFileName(*appointment.AfterImage))

			image = picture.FromID(&combinedImageId)
		}

		if image == nil {
			image = picture.FromID(appointment.AfterImage)
		}

		if image == nil {
			image = picture.FromID(appointment.BeforeImage)
		}

		reviewTitle += i18n.Language(ctx)["APPOINTMENT_TITLE"] + " " + i18n.FormatDate(ctx, appointment.Start)
	}

	if err != nil {
		w.WriteHeader(500)
		return
	}

	starsText := i18n.Language(ctx)["STARS_SINGULAR"]
	if review.Stars > 1 {
		starsText = i18n.Language(ctx)["STARS_PLURAL"]
	}

	templateParameters := map[string]string{
		"reviewTitle":    "",
		"reviewText":     fmt.Sprintf("%v %s - %s", review.Stars, starsText, review.Text),
		"linkUrl":        "",
		"branchImageUrl": "",
		"companyUrl":     share.ResolveCompanyUrlFromBranchId(ctx, h.Prisma, branch.ID),
		"siteName":       i18n.Language(ctx)["SITE_NAME"],
	}

	if branch != nil {
		companyName, _ := h.Prisma.Branch(prisma.BranchWhereUniqueInput{
			ID: &branch.ID,
		}).Company().Name().Exec(ctx)

		localizedCompanyName := i18n.GetLocalizedString(ctx, companyName)

		reviewTitle += " | "

		if localizedCompanyName != nil {
			reviewTitle += *localizedCompanyName + " "
		}

		branchName, _ := h.Prisma.Branch(prisma.BranchWhereUniqueInput{
			ID: &branch.ID,
		}).Name().Exec(ctx)

		localizedBranchName := i18n.GetLocalizedString(ctx, branchName)

		if localizedBranchName != nil {
			reviewTitle += *localizedBranchName
		}

		templateParameters["reviewTitle"] = reviewTitle

		if image == nil {
			image = picture.FromID(branch.Image)
		}

		templateParameters["linkUrl"] = share.ResolveShareRedirectUrl(ctx, h.Prisma, branch)
	}

	if image != nil {
		templateParameters["branchImageUrl"] = image.URL
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
