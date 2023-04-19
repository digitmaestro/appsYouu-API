package hours

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

var utc, _ = time.LoadLocation("UTC")

func FullDay(ctx context.Context, day prisma.DayOfWeek) string {
	return i18n.Language(ctx)["DAY_"+string(day)]
}

func ShortDay(ctx context.Context, day prisma.DayOfWeek) string {
	return i18n.Language(ctx)["DAY_KEY_"+string(day)]
}

func GetFirstAndLast(input string) string {
	fd := strings.Split(input, sep)
	return fd[0] + sep + fd[len(fd)-1]
}

func Timerange(loc *time.Location, start *string, end *string) *gqlgen.Timerange {
	if start == nil || end == nil {
		return nil
	}

	return &gqlgen.Timerange{
		Start: prisma.TimeDate(loc, *start),
		End:   prisma.TimeDate(loc, *end),
	}
}

func WeekdayToDayOfWeek(w time.Weekday) prisma.DayOfWeek {
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

func TodaySetTime(t *string) time.Time {
	start := time.Now()
	x := prisma.TimeDate(utc, *t)
	d := time.Date(
		start.Year(), start.Month(), start.Day(),
		x.Hour(), x.Minute(), x.Second(), x.Nanosecond(),
		utc,
	)
	return d
}
