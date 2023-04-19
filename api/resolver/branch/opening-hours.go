package branch

import (
	"context"
	"strings"
	"time"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/i18n"
	"github.com/steebchen/keskin-api/prisma"
)

var berlin, _ = time.LoadLocation("Europe/Berlin")

func (r *Branch) OpeningHours(ctx context.Context, obj *prisma.Branch) (*gqlgen.OpeningHours, error) {
	openingHours, err := r.Prisma.Branch(prisma.BranchWhereUniqueInput{
		ID: &obj.ID,
	}).OpeningHours(nil).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &gqlgen.OpeningHours{
		Formatted: formatOpeningHoursPtr(formatOpeningHours(ctx, berlin, openingHours)),
		Raw:       rawOpeningHoursPtr(rawOpeningHours(ctx, berlin, openingHours)),
	}, nil
}

func rawOpeningHoursPtr(input []gqlgen.RawOpeningHour) []*gqlgen.RawOpeningHour {
	nodes := []*gqlgen.RawOpeningHour{}

	for _, item := range input {
		clone := item
		nodes = append(nodes, &clone)
	}

	return nodes
}

func formatOpeningHoursPtr(input []gqlgen.FormattedOpeningHour) []*gqlgen.FormattedOpeningHour {
	nodes := []*gqlgen.FormattedOpeningHour{}

	for _, item := range input {
		clone := item
		nodes = append(nodes, &clone)
	}

	return nodes
}

const sep = "-"

// rawOpeningHours returns the raw database opening hours with all day pre-filled, regardless of whether they exist
// in the underlying database or not
func rawOpeningHours(ctx context.Context, loc *time.Location, openingHours []prisma.BranchOpeningHour) []gqlgen.RawOpeningHour {
	raw := []gqlgen.RawOpeningHour{}

	for _, day := range prisma.AllDayOfWeek {
		item := findByDay(openingHours, day)

		if item != nil {
			o := gqlgen.RawOpeningHour{
				Closed:   item.Closed,
				Break:    item.Break,
				Day:      item.Day,
				ShortDay: shortDay(ctx, day),
				FullDay:  fullDay(ctx, day),
			}

			if !item.Closed {
				o.Forenoon = timerange(loc, item.StartForenoon, item.EndForenoon)

				if item.Break {
					o.Afternoon = timerange(loc, item.StartAfternoon, item.EndAfternoon)
				}
			}

			raw = append(raw, o)
		} else {
			raw = append(raw, gqlgen.RawOpeningHour{
				Closed:   true,
				Break:    false,
				Day:      day,
				ShortDay: shortDay(ctx, day),
				FullDay:  fullDay(ctx, day),
			})
		}
	}

	return raw
}

func timerange(loc *time.Location, start *string, end *string) *gqlgen.Timerange {
	if start == nil || end == nil {
		return nil
	}

	return &gqlgen.Timerange{
		Start: prisma.TimeDate(loc, *start),
		End:   prisma.TimeDate(loc, *end),
	}
}

func formatOpeningHours(ctx context.Context, loc *time.Location, openingHours []prisma.BranchOpeningHour) []gqlgen.FormattedOpeningHour {
	formatted := []gqlgen.FormattedOpeningHour{}

	var previousValue string

	for _, day := range prisma.AllDayOfWeek {
		item := findByDay(openingHours, day)

		var value string
		var hasBreak bool
		var closed bool

		if item != nil && !item.Closed {
			if item.StartForenoon != nil && item.EndForenoon != nil {
				start := prisma.TimeDate(loc, *item.StartForenoon)
				end := prisma.TimeDate(loc, *item.EndForenoon)

				value += formatHour(ctx, start, end)
			}

			if item.Break && item.StartAfternoon != nil && item.EndAfternoon != nil {
				hasBreak = true

				start := prisma.TimeDate(loc, *item.StartAfternoon)
				end := prisma.TimeDate(loc, *item.EndAfternoon)

				value += " & "
				value += formatHour(ctx, start, end)
			}
		} else {
			closed = true
			value = i18n.Language(ctx)["CLOSED"]
		}

		if len(formatted) > 0 && previousValue == value {
			prev := &formatted[len(formatted)-1]

			// replace a single DayOfWeek with multiple,
			// e.g. "Mo" -> "Mo-Di"
			prev.Key = prev.Key + sep + string(day)
			prev.ShortDay = prev.ShortDay + sep + shortDay(ctx, day)
			prev.FullDay = prev.FullDay + sep + fullDay(ctx, day)

			// if more than two consecutive days are equal, cut off the days in between,
			// e.g. "Mo-Di-Mi-Do" -> "Mo-Do"
			prev.Key = getFirstAndLast(prev.Key)
			prev.ShortDay = getFirstAndLast(prev.ShortDay)
			prev.FullDay = getFirstAndLast(prev.FullDay)
		} else {
			formatted = append(formatted, gqlgen.FormattedOpeningHour{
				Key:      string(day),
				FullDay:  fullDay(ctx, day),
				ShortDay: shortDay(ctx, day),
				Value:    value,
				Closed:   closed,
				Break:    hasBreak,
			})
		}

		previousValue = value
	}

	return formatted
}

func formatHour(ctx context.Context, start time.Time, end time.Time) string {
	return i18n.FormatHourRange(ctx, start, end)
}

func fullDay(ctx context.Context, day prisma.DayOfWeek) string {
	return i18n.Language(ctx)["DAY_"+string(day)]
}

func shortDay(ctx context.Context, day prisma.DayOfWeek) string {
	return i18n.Language(ctx)["DAY_KEY_"+string(day)]
}

func findByDay(arr []prisma.BranchOpeningHour, key prisma.DayOfWeek) *prisma.BranchOpeningHour {
	for i, n := range arr {
		if key == n.Day {
			return &arr[i]
		}
	}
	return nil
}

func getFirstAndLast(input string) string {
	fd := strings.Split(input, sep)
	return fd[0] + sep + fd[len(fd)-1]
}
