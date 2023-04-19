package webmanifest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

type JSONManifestIcon struct {
	Src     string `json:"src"`
	Type    string `json:"type"`
	Sizes   string `json:"sizes"`
	Purpose string `json:"purpose"`
}

type JSONManifest struct {
	Name            string             `json:"name"`
	ShortName       string             `json:"short_name"`
	ThemeColor      string             `json:"theme_color"`
	BackgroundColor string             `json:"background_color"`
	Display         string             `json:"display"`
	Scope           string             `json:"scope"`
	StartURL        string             `json:"start_url"`
	GcmSenderId     string             `json:"gcm_sender_id"`
	Icons           []JSONManifestIcon `json:"icons"`
}

var GcmSenderId = ""

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	companyPwaShortName, err := h.Prisma.Company(prisma.CompanyWhereUniqueInput{
		ID: &companyId,
	}).PwaShortName().Exec(ctx)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	manifestCompanyName := strings.DefaultWhenEmpty(i18n.GetLocalizedString(ctx, companyName), "appsYouu")

	mainfest := JSONManifest{
		Name:            manifestCompanyName,
		ShortName:       strings.DefaultWhenEmpty(i18n.GetLocalizedString(ctx, companyPwaShortName), manifestCompanyName),
		ThemeColor:      strings.DefaultWhenEmpty(&company.PwaThemeColor, "#111111"),
		BackgroundColor: strings.DefaultWhenEmpty(&company.PwaBackgroundColor, "#111111"),
		Display:         "standalone",
		Scope:           "./",
		StartURL:        "./",
		GcmSenderId:     GcmSenderId,
	}

	if company.PwaIcon != nil {
		pwaIconID := picture.IDFromFileName(*company.PwaIcon)
		pwaIcon := picture.FromID(&pwaIconID)

		for _, size := range picture.IconSizes {
			mainfest.Icons = append(mainfest.Icons, JSONManifestIcon{
				Src:     fmt.Sprintf("%s_%v.png", pwaIcon.URL, size),
				Type:    "image/png",
				Sizes:   fmt.Sprintf("%vx%v", size, size),
				Purpose: "maskable any",
			})
		}
	} else {
		// defaults to icons included in Angular build
		for _, size := range []uint{72, 96, 128, 144, 152, 192, 384, 512} {
			mainfest.Icons = append(mainfest.Icons, JSONManifestIcon{
				Src:     fmt.Sprintf("assets/icons/icon-%vx%v.png", size, size),
				Type:    "image/png",
				Sizes:   fmt.Sprintf("%vx%v", size, size),
				Purpose: "maskable any",
			})
		}
	}

	manifestJson, err := json.Marshal(mainfest)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Write(manifestJson)
}

func New(client *prisma.Client) *Handler {
	return &Handler{
		Prisma: client,
	}
}
