package timeslots

import (
	"reflect"
	"testing"
	"time"

	"github.com/steebchen/keskin-api/prisma"
)

func Test_weekdayToDayOfWeek(t *testing.T) {
	type args struct {
		w time.Weekday
	}
	tests := []struct {
		name string
		args args
		want prisma.DayOfWeek
	}{
		{
			name: "mo",
			args: args{w: time.Monday},
			want: prisma.DayOfWeekMo,
		},
		{
			name: "tu",
			args: args{w: time.Tuesday},
			want: prisma.DayOfWeekTu,
		},
		{
			name: "we",
			args: args{w: time.Wednesday},
			want: prisma.DayOfWeekWe,
		},
		{
			name: "th",
			args: args{w: time.Thursday},
			want: prisma.DayOfWeekTh,
		},
		{
			name: "fr",
			args: args{w: time.Friday},
			want: prisma.DayOfWeekFr,
		},
		{
			name: "sa",
			args: args{w: time.Saturday},
			want: prisma.DayOfWeekSa,
		},
		{
			name: "su",
			args: args{w: time.Sunday},
			want: prisma.DayOfWeekSu,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := weekdayToDayOfWeek(tt.args.w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("weekdayToDayOfWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatTime(t *testing.T) {
	type args struct {
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "a",
			args: args{
				date: prisma.TimeDate(utc, "2000-01-05T07:35:59.999Z"),
			},
			want: "07:35",
		},
		{
			name: "b",
			args: args{
				date: prisma.TimeDate(utc, "2000-01-01T00:00:00.000Z"),
			},
			want: "00:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatTime(tt.args.date); got != tt.want {
				t.Errorf("formatTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
