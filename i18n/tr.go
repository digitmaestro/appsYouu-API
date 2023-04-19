package i18n

import (
	"fmt"
	"time"
)

var TR = map[string]string{
	"DAY_MO": "Pazartesi",
	"DAY_TU": "Sali",
	"DAY_WE": "Çarşamba",
	"DAY_TH": "Perşembe",
	"DAY_FR": "Cuma",
	"DAY_SA": "Cumartesi",
	"DAY_SU": "Pazar",

	"DAY_KEY_MO": "Pzt",
	"DAY_KEY_TU": "Sal",
	"DAY_KEY_WE": "Çar",
	"DAY_KEY_TH": "Per",
	"DAY_KEY_FR": "Cum",
	"DAY_KEY_SA": "Cmt",
	"DAY_KEY_SU": "Paz",

	"TIME":        "saatleri arası",
	"CLOSED":      "kapalı",
	"NOT_WORKING": "çalışmıyor",

	"APPOINTMENT_REMINDER_TITLE":         "Yarınki randevunuz",
	"APPOINTMENT_REMINDER_TEXT":          "Bizimle %s tarihinde saat %s de olan randevunuzu size hatırlatmak istiyoruz.",
	"APPOINTMENT_APPROVED_TITLE":         "Bizimle bir sonraki randevunuz",
	"APPOINTMENT_APPROVED_TEXT":          "Bizimle olan bir sonraki randevunuzu %s, %s de onaylıyoruz.",
	"APPOINTMENT_CANCELED_TITLE":         "Randevu iptali",
	"APPOINTMENT_CANCELED_TEXT":          "Üzgünüzki %s tarihinde saat %s de olan randevunuzu  de iptal etmemiz gerektiğini bildiririz.",
	"APPOINTMENT_CANCELED_BY_USER_TITLE": "Randevu iptali",
	"APPOINTMENT_CANCELED_BY_USER_TEXT":  "Randevunuz iptal edilmişdir.",
	"BIRTHDAY_TITLE":                     "Doğum gününüz kutlu olsun",
	"BIRTHDAY_TEXT":                      "Doğum gününüzü kutluyor ve bir sonraki ziyaretinizi dört gözle bekliyoruz.",
	"TEST_NOTIFICATION_TITLE":            "Deney bildirimi",
	"TEST_NOTIFICATION_TEXT":             "Bunu okuyabiliyorsanız, her şeyi doğru yapmışsınızdır.",

	"STARS_SINGULAR":    "Yıldız",
	"STARS_PLURAL":      "Yıldızlar",
	"APPOINTMENT_TITLE": "Randevu tarihi",

	"SALUTATION_MALE":   "Bey",
	"SALUTATION_FEMALE": "Hanım",

	"SITE_NAME": "appsYouu",

	"DELETED_PRODUCT": "Silinmiş ürün",
	"DELETED_SERVICE": "Silinmiş hizmet",
	"DELETED_USER":    "Silinmiş kullanıcı",
}

func formatDateTR(t time.Time) string {
	return t.Format("02.01.2006")
}

func formatTimeTR(t time.Time) string {
	f := fmt.Sprintf("%d:%d", t.Hour(), t.Minute())

	if t.Minute() == 0 {
		f = fmt.Sprintf("%d", t.Hour())
	}

	return f
}

func formatHourRangeTR(start time.Time, end time.Time) string {
	s := formatTimeTR(start)
	e := formatTimeTR(end)

	return fmt.Sprintf("%s-%s %s", s, e, TR["TIME"])
}
