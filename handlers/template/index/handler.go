package index

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/cbroglie/mustache"

	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/lib/auth"
	"github.com/steebchen/keskin-api/lib/picture"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/lib/strings"
	"github.com/steebchen/keskin-api/prisma"
)

type Handler struct {
	Prisma *prisma.Client
}

var TemplateFolder = "/templates/"

// included in Angular build
const defaultAppleTouchIconsTemplate = `
<link rel="apple-touch-icon" href="assets/icons/icon-512x512.png">
<link rel="apple-touch-icon" sizes="152x152" href="assets/icons/icon-152x152.png">
<link rel="apple-touch-icon" sizes="192x192" href="assets/icons/icon-192x192.png">
`

// included in Angular build
const defaultFaviconsTemplate = `
<link href="assets/icon/favicon.png" rel="icon" type="image/png"/>
`

const appleTouchIconsTemplate = `
<link rel="apple-touch-icon" href="{{iconUrl}}_512.png">
<link rel="apple-touch-icon" sizes="152x152" href="{{iconUrl}}_152.png">
<link rel="apple-touch-icon" sizes="192x192" href="{{iconUrl}}_192.png">
`

const faviconsTemplate = `
<link rel="icon" type="image/png" sizes="64x64" href="{{iconUrl}}_64.png">
<link rel="icon" type="image/png" sizes="32x32" href="{{iconUrl}}_32.png">
<link rel="icon" type="image/png" sizes="16x16" href="{{iconUrl}}_16.png">
`

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	htmlTemplate, err := ioutil.ReadFile(filepath.Join(TemplateFolder, "index.html"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	ctx := r.Context()
	ctx = sessctx.SetCompanyHeader(ctx, r.Header.Get(auth.CompanyHeader))
	ctx = sessctx.SetHost(ctx, r.Host)
	language := r.Header.Get(auth.LanguageHeader)
	ctx = sessctx.SetLanguage(ctx, &language)

	companyId := sessctx.CompanyWithFallback(ctx, h.Prisma, nil)

	company, err := h.Prisma.Company(prisma.CompanyWhereUniqueInput{
		ID: &companyId,
	}).Exec(ctx)

	if err != nil || company == nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}

	companyName, err := h.Prisma.Company(prisma.CompanyWhereUniqueInput{
		ID: &companyId,
	}).Name().Exec(ctx)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	companyNameString := i18n.GetLocalizedString(ctx, companyName)

	templateParameters := map[string]string{
		"title":      strings.DefaultWhenEmpty(companyNameString, "appsYouu"),
		"themeColor": strings.DefaultWhenEmpty(&company.PwaThemeColor, "#111111"),
	}

	if company.PwaIcon == nil {
		templateParameters["appleTouchIcons"] = defaultAppleTouchIconsTemplate
		templateParameters["favicons"] = defaultFaviconsTemplate
	} else {
		pwaIconID := picture.IDFromFileName(*company.PwaIcon)
		pwaIcon := picture.FromID(&pwaIconID)

		filledAppleTouchIconsTemplate, err := mustache.Render(appleTouchIconsTemplate, map[string]string{
			"iconUrl": pwaIcon.URL,
		})

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		templateParameters["appleTouchIcons"] = filledAppleTouchIconsTemplate

		filledFaviconsTemplate, err := mustache.Render(faviconsTemplate, map[string]string{
			"iconUrl": pwaIcon.URL,
		})

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}

		templateParameters["favicons"] = filledFaviconsTemplate
	}

	filledTemplate, err := mustache.Render(string(htmlTemplate), templateParameters)

	if err != nil {
		log.Println(err)
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
