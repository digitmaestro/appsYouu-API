package i18n

import (
	"fmt"
	"time"
)

var EN = map[string]string{
	"DAY_MO": "Monday",
	"DAY_TU": "Tuesday",
	"DAY_WE": "Wednesday",
	"DAY_TH": "Thursday",
	"DAY_FR": "Friday",
	"DAY_SA": "Saturday",
	"DAY_SU": "Sunday",

	"DAY_KEY_MO": "Mon",
	"DAY_KEY_TU": "Tue",
	"DAY_KEY_WE": "Wed",
	"DAY_KEY_TH": "Thu",
	"DAY_KEY_FR": "Fri",
	"DAY_KEY_SA": "Sat",
	"DAY_KEY_SU": "Sun",

	"CLOSED":      "closed",
	"NOT_WORKING": "not working",

	"APPOINTMENT_REMINDER_TITLE":         "Your appoinment with us tomorrow",
	"APPOINTMENT_REMINDER_TEXT":          "We'd like to remind you about our appointment on %s at %s.",
	"APPOINTMENT_APPROVED_TITLE":         "Your next appoinment with us",
	"APPOINTMENT_APPROVED_TEXT":          "We'd like to confirm your appointment on %s at %s.",
	"APPOINTMENT_CANCELED_TITLE":         "Cacelation",
	"APPOINTMENT_CANCELED_TEXT":          "We regret to inform you that we need to cancel your appointment with us on %s at %s.",
	"APPOINTMENT_CANCELED_BY_USER_TITLE": "Cacelation",
	"APPOINTMENT_CANCELED_BY_USER_TEXT":  "We have received your appointment cancellation.",
	"BIRTHDAY_TITLE":                     "Happy Birthday",
	"BIRTHDAY_TEXT":                      "we wish you a happy birthday and are looking forward to see you again soon.",
	"TEST_NOTIFICATION_TITLE":            "Test notification",
	"TEST_NOTIFICATION_TEXT":             "When you can read this, you everything is working well.",

	"STARS_SINGULAR":    "Star",
	"STARS_PLURAL":      "Stars",
	"APPOINTMENT_TITLE": "Appointment at",

	"SALUTATION_MALE":   "Hello",
	"SALUTATION_FEMALE": "Hello",

	"SITE_NAME": "appsYouu",

	"DELETED_PRODUCT": "Deleted Product",
	"DELETED_SERVICE": "Deleted Service",
	"DELETED_USER":    "Deleted User",
}

func formatDateEN(t time.Time) string {
	return t.Format("02/01/2006")
}

func formatTimeEN(t time.Time) string {
	f := t.Format("3:04pm")

	if t.Minute() == 0 {
		f = t.Format("3pm")
	}

	return f
}

func formatHourRangeEN(start time.Time, end time.Time) string {
	s := formatTimeEN(start)
	e := formatTimeEN(end)

	return fmt.Sprintf("%s-%s", s, e)
}
