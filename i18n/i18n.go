package i18n

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

var DefaultLanguage = "DE"
var berlin, _ = time.LoadLocation("Europe/Berlin")

func IsAllowedLanguage(language string) bool {
	return language == "DE" || language == "EN" || language == "TR"
}

func languageOrDefault(language *string) string {
	if language != nil && IsAllowedLanguage(*language) {
		return *language
	} else {
		return DefaultLanguage
	}
}

func SessionLanguageKey(ctx context.Context) string {
	return languageOrDefault(sessctx.Language(ctx))
}

func Language(ctx context.Context) map[string]string {
	switch SessionLanguageKey(ctx) {
	case "EN":
		return EN
	case "TR":
		return TR
	default:
		return DE
	}
}

func toLocalDate(dateString string) time.Time {
	localDate, err := time.Parse(time.RFC3339, dateString)

	if err != nil {
		panic(err)
	}

	return localDate.In(berlin)
}

func FormatDate(ctx context.Context, utcDate string) string {
	localDate := toLocalDate(utcDate)

	switch SessionLanguageKey(ctx) {
	case "EN":
		return formatDateEN(localDate)
	case "TR":
		return formatDateTR(localDate)
	default:
		return formatDateDE(localDate)
	}
}

func FormatTime(ctx context.Context, utcDate string) string {
	localDate := toLocalDate(utcDate)

	switch SessionLanguageKey(ctx) {
	case "EN":
		return formatTimeEN(localDate)
	case "TR":
		return formatTimeTR(localDate)
	default:
		return formatTimeDE(localDate)
	}
}

func FormatHourRange(ctx context.Context, start time.Time, end time.Time) string {
	switch SessionLanguageKey(ctx) {
	case "EN":
		return formatHourRangeEN(start, end)
	case "TR":
		return formatHourRangeTR(start, end)
	default:
		return formatHourRangeDE(start, end)
	}
}

func FormatPrice(ctx context.Context, price float64) string {
	switch SessionLanguageKey(ctx) {
	case "EN":
		return fmt.Sprintf("%.2f", price)
	default:
		return strings.Replace(fmt.Sprintf("%.2f", price), ".", ",", 1)
	}

	return ""
}

func CreateLocalizedString(ctx context.Context, value *string) *prisma.LocalizedStringCreateOneInput {
	create := prisma.LocalizedStringCreateInput{}

	switch SessionLanguageKey(ctx) {
	case "EN":
		create.En = value
	case "TR":
		create.Tr = value
	default:
		create.De = value
	}

	return &prisma.LocalizedStringCreateOneInput{
		Create: &create,
	}
}

func LocalizedStringUpdateDataInput(ctx context.Context, value *string) *prisma.LocalizedStringUpdateDataInput {
	update := prisma.LocalizedStringUpdateDataInput{}

	switch SessionLanguageKey(ctx) {
	case "EN":
		update.En = value
	case "TR":
		update.Tr = value
	default:
		update.De = value
	}

	return &update
}

func UpdateRequiredLocalizedString(ctx context.Context, value *string) *prisma.LocalizedStringUpdateOneRequiredInput {
	return &prisma.LocalizedStringUpdateOneRequiredInput{
		Update: LocalizedStringUpdateDataInput(ctx, value),
	}
}

func GetLocalizedString(ctx context.Context, localizations *prisma.LocalizedString) *string {
	var result *string = nil

	if localizations != nil {
		switch SessionLanguageKey(ctx) {
		case "EN":
			result = localizations.En
		case "TR":
			result = localizations.Tr
		default:
			result = localizations.De
		}

		if result == nil {
			switch DefaultLanguage {
			case "EN":
				result = localizations.En
			case "TR":
				result = localizations.Tr
			default:
				result = localizations.De
			}
		}
	}

	return result
}
