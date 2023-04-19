package timeslots

import (
	"fmt"
	"time"

	"github.com/steebchen/keskin-api/prisma"
)

func weekdayToDayOfWeek(w time.Weekday) prisma.DayOfWeek {
	switch w {
	case time.Monday:
		return prisma.DayOfWeekMo
	case time.Tuesday:
		return prisma.DayOfWeekTu
	case time.Wednesday:
		return prisma.DayOfWeekWe
	case time.Thursday:
		return prisma.DayOfWeekTh
	case time.Friday:
		return prisma.DayOfWeekFr
	case time.Saturday:
		return prisma.DayOfWeekSa
	case time.Sunday:
		return prisma.DayOfWeekSu
	default:
		panic(fmt.Sprintf("could not convert weekday %d", w))
	}
}

func formatTime(date time.Time) string {
	return fmt.Sprintf("%02d:%02d", date.Hour(), date.Minute())
}
