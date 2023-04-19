package hours

import (
	"context"
	"time"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

var berlin, _ = time.LoadLocation("Europe/Berlin")

func RawWorkingHoursPtr(input []gqlgen.RawWorkingHours) []*gqlgen.RawWorkingHours {
	nodes := []*gqlgen.RawWorkingHours{}

	for _, item := range input {
		clone := item
		nodes = append(nodes, &clone)
	}

	return nodes
}

func FormatWorkingHoursPtr(input []gqlgen.FormattedWorkingHours) []*gqlgen.FormattedWorkingHours {
	nodes := []*gqlgen.FormattedWorkingHours{}

	for _, item := range input {
		clone := item
		nodes = append(nodes, &clone)
	}

	return nodes
}

const sep = "-"

// RawWorkingHours returns the raw database working hours with all day pre-filled, regardless of whether they exist
// in the underlying database or not
func RawWorkingHours(ctx context.Context, loc *time.Location, workingHours []prisma.WorkingHours) []gqlgen.RawWorkingHours {
	raw := []gqlgen.RawWorkingHours{}

	for _, day := range prisma.AllDayOfWeek {
		item := findByDay(workingHours, day)

		if item != nil {
			o := gqlgen.RawWorkingHours{
				NotWorking: item.NotWorking,
				Break:      item.Break,
				Day:        item.Day,
				ShortDay:   ShortDay(ctx, day),
				FullDay:    FullDay(ctx, day),
			}

			if !item.NotWorking {
				o.Forenoon = Timerange(loc, item.StartForenoon, item.EndForenoon)

				if item.Break {
					o.Afternoon = Timerange(loc, item.StartAfternoon, item.EndAfternoon)
				}
			}

			raw = append(raw, o)
		} else {
			raw = append(raw, gqlgen.RawWorkingHours{
				NotWorking: true,
				Break:      false,
				Day:        day,
				ShortDay:   ShortDay(ctx, day),
				FullDay:    FullDay(ctx, day),
			})
		}
	}

	return raw
}

func FormatWorkingHours(ctx context.Context, loc *time.Location, workingHours []prisma.WorkingHours) []gqlgen.FormattedWorkingHours {
	formatted := []gqlgen.FormattedWorkingHours{}

	var previousValue string

	for _, day := range prisma.AllDayOfWeek {
		item := findByDay(workingHours, day)

		var value string
		var hasBreak bool
		var notWorking bool

		if item != nil && !item.NotWorking {
			if item.StartForenoon != nil && item.EndForenoon != nil {
				start := prisma.TimeDate(loc, *item.StartForenoon)
				end := prisma.TimeDate(loc, *item.EndForenoon)

				value += i18n.FormatHourRange(ctx, start, end)
			}

			if item.Break && item.StartAfternoon != nil && item.EndAfternoon != nil {
				hasBreak = true

				start := prisma.TimeDate(loc, *item.StartAfternoon)
				end := prisma.TimeDate(loc, *item.EndAfternoon)

				value += " & "
				value += i18n.FormatHourRange(ctx, start, end)
			}
		} else {
			notWorking = true
			value = i18n.Language(ctx)["NOT_WORKING"]
		}

		if len(formatted) > 0 && previousValue == value {
			prev := &formatted[len(formatted)-1]

			// replace a single DayOfWeek with multiple,
			// e.g. "Mo" -> "Mo-Di"
			prev.Key = prev.Key + sep + string(day)
			prev.ShortDay = prev.ShortDay + sep + ShortDay(ctx, day)
			prev.FullDay = prev.FullDay + sep + FullDay(ctx, day)

			// if more than two consecutive days are equal, cut off the days in between,
			// e.g. "Mo-Di-Mi-Do" -> "Mo-Do"
			prev.Key = GetFirstAndLast(prev.Key)
			prev.ShortDay = GetFirstAndLast(prev.ShortDay)
			prev.FullDay = GetFirstAndLast(prev.FullDay)
		} else {
			formatted = append(formatted, gqlgen.FormattedWorkingHours{
				Key:        string(day),
				FullDay:    FullDay(ctx, day),
				ShortDay:   ShortDay(ctx, day),
				Value:      value,
				NotWorking: notWorking,
				Break:      hasBreak,
			})
		}

		previousValue = value
	}

	return formatted
}

func findByDay(arr []prisma.WorkingHours, key prisma.DayOfWeek) *prisma.WorkingHours {
	for i, n := range arr {
		if key == n.Day {
			return &arr[i]
		}
	}
	return nil
}
