package i18n

import (
	"fmt"
	"time"
)

var DE = map[string]string{
	"DAY_MO": "Montag",
	"DAY_TU": "Dienstag",
	"DAY_WE": "Mittwoch",
	"DAY_TH": "Donnerstag",
	"DAY_FR": "Freitag",
	"DAY_SA": "Samstag",
	"DAY_SU": "Sonntag",

	"DAY_KEY_MO": "Mo",
	"DAY_KEY_TU": "Di",
	"DAY_KEY_WE": "Mi",
	"DAY_KEY_TH": "Do",
	"DAY_KEY_FR": "Fr",
	"DAY_KEY_SA": "Sa",
	"DAY_KEY_SU": "So",

	"TIME":        "Uhr",
	"CLOSED":      "geschlossen",
	"NOT_WORKING": "arbeitet nicht",

	"APPOINTMENT_REMINDER_TITLE":         "Ihr Termin morgen bei uns",
	"APPOINTMENT_REMINDER_TEXT":          "Wir möchten Sie gerne an Ihren Termin am %s um %s Uhr bei uns erinnern.",
	"APPOINTMENT_APPROVED_TITLE":         "Ihr nächster Termin bei uns",
	"APPOINTMENT_APPROVED_TEXT":          "Hiermit bestätigen wir Ihren nächsten Termin bei uns am %s um %s Uhr.",
	"APPOINTMENT_CANCELED_TITLE":         "Terminabsage",
	"APPOINTMENT_CANCELED_TEXT":          "Mit Bedauern teilen wir Ihnen mit, dass wir Ihren Termin am %s von %s Uhr bei uns absagen müssen.",
	"APPOINTMENT_CANCELED_BY_USER_TITLE": "Terminabsage",
	"APPOINTMENT_CANCELED_BY_USER_TEXT":  "Ihre Terminabsage ist bei uns eingegangen.",
	"BIRTHDAY_TITLE":                     "Herzlichen Glückwunsch zum Geburtstag",
	"BIRTHDAY_TEXT":                      "Wir wünschen Ihnen alles Gute zum Geburtstag und freuen uns auf Ihren nächsten Besuch.",
	"TEST_NOTIFICATION_TITLE":            "Testbenachrichtigung",
	"TEST_NOTIFICATION_TEXT":             "Wenn Sie das lesen können, dann haben Sie alles richtig gemacht.",

	"STARS_SINGULAR":    "Stern",
	"STARS_PLURAL":      "Sterne",
	"APPOINTMENT_TITLE": "Termin am",

	"SALUTATION_MALE":   "Sehr geehrter Herr",
	"SALUTATION_FEMALE": "Sehr geehrte Frau",

	"SITE_NAME": "appsYouu",

	"DELETED_PRODUCT": "Gelöschtes Produkt",
	"DELETED_SERVICE": "Gelöschte Dienstleistung",
	"DELETED_USER":    "Gelöschter Nutzer",
}

func formatDateDE(t time.Time) string {
	return t.Format("02.01.2006")
}

func formatTimeDE(t time.Time) string {
	f := fmt.Sprintf("%d:%d", t.Hour(), t.Minute())

	if t.Minute() == 0 {
		f = fmt.Sprintf("%d", t.Hour())
	}

	return f
}

func formatHourRangeDE(start time.Time, end time.Time) string {
	s := formatTimeDE(start)
	e := formatTimeDE(end)

	return fmt.Sprintf("%s-%s %s", s, e, DE["TIME"])
}
