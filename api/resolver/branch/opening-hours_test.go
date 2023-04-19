package branch

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/lib/sessctx"
	"github.com/steebchen/keskin-api/prisma"
)

var utc, _ = time.LoadLocation("UTC")

func hour(hour time.Duration, minute time.Duration) string {
	return time.Now().
		UTC().
		Truncate(time.Hour * 24).
		Add(hour * time.Hour).
		Add(minute * time.Minute).
		Format(time.RFC3339)
}

func Test_hour(t *testing.T) {
	h := hour(5, 20)
	got, _ := time.Parse(time.RFC3339, h)
	assert.Equal(t, 5, got.UTC().Hour())
	assert.Equal(t, 20, got.UTC().Minute())
}

func Test_formatOpeningHours(t *testing.T) {
	type args struct {
		openingHours []prisma.BranchOpeningHour
	}
	tests := []struct {
		name string
		args args
		want []gqlgen.FormattedOpeningHour
	}{
		{
			name: "single full day",
			args: args{
				openingHours: []prisma.BranchOpeningHour{
					{
						Day: prisma.DayOfWeekMo,
						// opens 9-5
						StartForenoon: prisma.Str(hour(9, 0)),
						EndForenoon:   prisma.Str(hour(17, 0)),
					},
				},
			},
			want: []gqlgen.FormattedOpeningHour{
				{
					Key:      "MO",
					ShortDay: "Mo",
					FullDay:  "Montag",
					Value:    "9-17 Uhr",
					Closed:   false,
					Break:    false,
				},
				{
					Key:      "TU-SU",
					ShortDay: "Di-So",
					FullDay:  "Dienstag-Sonntag",
					Value:    "geschlossen",
					Closed:   true,
					Break:    false,
				},
			},
		},
		{
			name: "forenoon and afternoon",
			args: args{
				openingHours: []prisma.BranchOpeningHour{
					{
						Day:            prisma.DayOfWeekMo,
						Break:          true,
						StartForenoon:  prisma.Str(hour(8, 0)),
						EndForenoon:    prisma.Str(hour(12, 0)),
						StartAfternoon: prisma.Str(hour(13, 0)),
						EndAfternoon:   prisma.Str(hour(18, 0)),
					},
				},
			},
			want: []gqlgen.FormattedOpeningHour{
				{
					Key:      "MO",
					ShortDay: "Mo",
					FullDay:  "Montag",
					Value:    "8-12 Uhr & 13-18 Uhr",
					Closed:   false,
					Break:    true,
				},
				{
					Key:      "TU-SU",
					ShortDay: "Di-So",
					FullDay:  "Dienstag-Sonntag",
					Value:    "geschlossen",
					Closed:   true,
					Break:    false,
				},
			},
		},
		{
			name: "minutes",
			args: args{
				openingHours: []prisma.BranchOpeningHour{
					{
						Day:           prisma.DayOfWeekMo,
						StartForenoon: prisma.Str(hour(8, 15)),
						EndForenoon:   prisma.Str(hour(17, 30)),
					},
				},
			},
			want: []gqlgen.FormattedOpeningHour{
				{
					Key:      "MO",
					ShortDay: "Mo",
					FullDay:  "Montag",
					Value:    "8:15-17:30 Uhr",
					Closed:   false,
					Break:    false,
				},
				{
					Key:      "TU-SU",
					ShortDay: "Di-So",
					FullDay:  "Dienstag-Sonntag",
					Value:    "geschlossen",
					Closed:   true,
					Break:    false,
				},
			},
		},
		{
			name: "multiple days",
			args: args{
				openingHours: []prisma.BranchOpeningHour{
					{
						Day:           prisma.DayOfWeekMo,
						StartForenoon: prisma.Str(hour(8, 0)),
						EndForenoon:   prisma.Str(hour(17, 0)),
					},
					{
						Day:           prisma.DayOfWeekTh,
						StartForenoon: prisma.Str(hour(9, 0)),
						EndForenoon:   prisma.Str(hour(15, 0)),
					},
				},
			},
			want: []gqlgen.FormattedOpeningHour{
				{
					Key:      "MO",
					ShortDay: "Mo",
					FullDay:  "Montag",
					Value:    "8-17 Uhr",
					Closed:   false,
					Break:    false,
				},
				{
					Key:      "TU-WE",
					ShortDay: "Di-Mi",
					FullDay:  "Dienstag-Mittwoch",
					Value:    "geschlossen",
					Closed:   true,
					Break:    false,
				},
				{
					Key:      "TH",
					ShortDay: "Do",
					FullDay:  "Donnerstag",
					Value:    "9-15 Uhr",
					Closed:   false,
					Break:    false,
				},
				{
					Key:      "FR-SU",
					ShortDay: "Fr-So",
					FullDay:  "Freitag-Sonntag",
					Value:    "geschlossen",
					Closed:   true,
					Break:    false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			language := "DE"
			ctx = sessctx.SetLanguage(ctx, &language)

			got := formatOpeningHours(ctx, utc, tt.args.openingHours)
			assert.Equal(t, tt.want, got)
		})
	}
}
