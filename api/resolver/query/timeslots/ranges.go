package timeslots

import (
	"context"
	"sort"
	"time"

	"github.com/steebchen/keskin-api/gqlgen"
	"github.com/steebchen/keskin-api/prisma"
)

var defaultInterval = 5 * time.Minute
var utc, _ = time.LoadLocation("UTC")

// findOpeningHour returns the matching openingHour from the current time
// if none is found, it returns a default openingHour using false default values which equals to being closed
func findOpeningHour(
	date *time.Time,
	openingHours []prisma.BranchOpeningHour,
) prisma.BranchOpeningHour {
	if date == nil {
		now := time.Now()
		date = &now
	}

	currentWeekday := weekdayToDayOfWeek(date.Weekday())
	var openingHour prisma.BranchOpeningHour

	for _, i := range openingHours {
		if i.Day == currentWeekday {
			openingHour = i
			break
		}
	}

	return openingHour
}

// findWorkingHuors returns the matching workingHours from the current time
// if none is found, it returns a default workingHours using false default values which equals to not working
func findWorkingHours(
	date *time.Time,
	workingHours []prisma.WorkingHours,
) []prisma.WorkingHours {
	if date == nil {
		now := time.Now()
		date = &now
	}

	currentWeekday := weekdayToDayOfWeek(date.Weekday())

	result := []prisma.WorkingHours{}

	for _, i := range workingHours {
		if i.Day == currentWeekday && !i.NotWorking {
			result = append(result, i)
		}
	}

	return result
}

func todaySetTime(start time.Time, t *string) time.Time {
	x := prisma.TimeDate(utc, *t)
	d := time.Date(
		start.Year(), start.Month(), start.Day(),
		x.Hour(), x.Minute(), x.Second(), x.Nanosecond(),
		utc,
	)
	return d
}

type TimeSlice struct {
	Start time.Time
	End   time.Time
}

func beforeOrEqual(t time.Time, u time.Time) bool {
	return t.Before(u) || t.Equal(u)
}

func afterOrEqual(t time.Time, u time.Time) bool {
	return t.After(u) || t.Equal(u)
}

// 1.
// slice    |---|
// newSlice       |---|
// result   |---| |---|
//
// 2.
// slice          |---|
// newSlice |---|
// result   |---| |---|
//
// 3.
// slice    |---|
// newSlice   |---|
// result   |-----|
//
// 4.
// slice      |---|
// newSlice |---|
// result   |-----|
//
// 5.
// slice    |-------|
// newSlice   |---|
// result   |-------|
//
// 6.
// slice      |---|
// newSlice |-------|
// result   |-------|

func join(slices []*TimeSlice, newSlice *TimeSlice) []*TimeSlice {
	if newSlice == nil {
		return slices
	}

	merged := false
	dropped := false

	for _, slice := range slices {
		if !merged {
			if beforeOrEqual(slice.Start, newSlice.Start) && afterOrEqual(slice.End, newSlice.End) {
				// 5.
				dropped = true
			} else if afterOrEqual(slice.Start, newSlice.Start) && beforeOrEqual(slice.End, newSlice.End) {
				// 6.
				slice.Start = newSlice.Start
				slice.End = newSlice.End
				merged = true
			} else if beforeOrEqual(slice.Start, newSlice.Start) && afterOrEqual(slice.End, newSlice.Start) && beforeOrEqual(slice.End, newSlice.End) {
				// 3.
				slice.End = newSlice.End
				merged = true
			} else if beforeOrEqual(slice.Start, newSlice.End) && afterOrEqual(slice.End, newSlice.End) && afterOrEqual(slice.Start, newSlice.Start) {
				// 4.
				slice.Start = newSlice.Start
				merged = true
			}
		}
	}

	if !dropped {
		if merged {
			for i := 0; i < len(slices)-1; {
				slice := slices[i]
				nextSlice := slices[i+1]

				if slice.End.After(nextSlice.Start) {
					if slice.End.Before(nextSlice.End) {
						slice.End = nextSlice.End
					}

					copy(slices[i+1:], slices[i+2:])
					slices = slices[:len(slices)-1]
				} else {
					i++
				}
			}
		} else {
			// 1., 2.
			slices = append(slices, newSlice)

			sort.Slice(slices, func(i, j int) bool {
				return slices[i].Start.Before(slices[j].Start)
			})
		}
	}

	return slices
}

// 1.
// slice     |---|
// intersect       |---|
// result
//
// 2.
// slice           |---|
// intersect |---|
// result
//
// 3.
// slice     |---|
// intersect   |---|
// result      |-|
//
// 4.
// slice       |---|
// intersect |---|
// result      |-|
//
// 5.
// slice     |-------|
// intersect   |---|
// result      |---|
//
// 6.
// slice       |---|
// intersect |-------|
// result      |---|

func intersect(slice *TimeSlice, intersectSlice *TimeSlice) *TimeSlice {
	if intersectSlice == nil || slice == nil || beforeOrEqual(slice.End, intersectSlice.Start) || afterOrEqual(slice.Start, intersectSlice.End) {
		// 1., 2.
		return nil
	} else if beforeOrEqual(slice.Start, intersectSlice.Start) && afterOrEqual(slice.End, intersectSlice.Start) && beforeOrEqual(slice.End, intersectSlice.End) {
		// 3.
		slice.Start = intersectSlice.Start
	} else if afterOrEqual(slice.Start, intersectSlice.Start) && beforeOrEqual(slice.Start, intersectSlice.End) && afterOrEqual(slice.End, intersectSlice.End) {
		// 4.
		slice.End = intersectSlice.End
	} else if beforeOrEqual(slice.Start, intersectSlice.Start) && afterOrEqual(slice.End, intersectSlice.End) {
		// 5.
		slice.Start = intersectSlice.Start
		slice.End = intersectSlice.End
	}
	// 6.

	return slice
}

func calculateWorkingHoursDuringOpeningHours(
	start time.Time,
	openingHour prisma.BranchOpeningHour,
	workingHoursSlices []*TimeSlice,
) []*TimeSlice {
	availableTimeSlices := []*TimeSlice{}

	var openingHoursForenoonTimeSlice *TimeSlice = nil
	if !openingHour.Closed &&
		openingHour.StartForenoon != nil &&
		openingHour.EndForenoon != nil {
		openingHoursForenoonTimeSlice = &TimeSlice{
			Start: todaySetTime(start, openingHour.StartForenoon),
			End:   todaySetTime(start, openingHour.EndForenoon),
		}
	}

	var openingHoursAfternoonTimeSlice *TimeSlice = nil
	if !openingHour.Closed &&
		openingHour.Break &&
		openingHour.StartAfternoon != nil &&
		openingHour.EndAfternoon != nil {
		openingHoursAfternoonTimeSlice = &TimeSlice{
			Start: todaySetTime(start, openingHour.StartAfternoon),
			End:   todaySetTime(start, openingHour.EndAfternoon),
		}
	}

	for _, workingHoursSlice := range workingHoursSlices {
		availableTimeSlices = join(
			availableTimeSlices,
			intersect(&TimeSlice{
				Start: workingHoursSlice.Start,
				End:   workingHoursSlice.End,
			}, openingHoursForenoonTimeSlice),
		)

		availableTimeSlices = join(
			availableTimeSlices,
			intersect(&TimeSlice{
				Start: workingHoursSlice.Start,
				End:   workingHoursSlice.End,
			}, openingHoursAfternoonTimeSlice),
		)
	}

	return availableTimeSlices
}

func calculateWorkingHoursSlices(
	prismaClient *prisma.Client,
	start time.Time,
	appointments []prisma.Appointment,
	workingHours []prisma.WorkingHours,
) []*TimeSlice {
	if len(workingHours) == 0 {
		return []*TimeSlice{}
	}

	ctx := context.Background()

	workingHoursByUser := make(map[string][]prisma.WorkingHours)

	for _, workingHoursItem := range workingHours {
		user, err := prismaClient.WorkingHours(prisma.WorkingHoursWhereUniqueInput{
			ID: &workingHoursItem.ID,
		}).User().Exec(ctx)

		if err == nil {
			_, exists := workingHoursByUser[user.ID]
			if !exists {
				workingHoursByUser[user.ID] = []prisma.WorkingHours{}
			}

			workingHoursByUser[user.ID] = append(workingHoursByUser[user.ID], workingHoursItem)
		}
	}

	appointmentsByUser := make(map[string][]prisma.Appointment)

	for _, appointment := range appointments {
		user, err := prismaClient.Appointment(prisma.AppointmentWhereUniqueInput{
			ID: &appointment.ID,
		}).Employee().Exec(ctx)

		if err == nil {
			_, exists := appointmentsByUser[user.ID]
			if !exists {
				appointmentsByUser[user.ID] = []prisma.Appointment{}
			}

			appointmentsByUser[user.ID] = append(appointmentsByUser[user.ID], appointment)
		}
	}

	workingHoursSlices := []*TimeSlice{}

	for user, userWorkingHours := range workingHoursByUser {
		availableStart := []time.Time{}
		availableEnd := []time.Time{}

		for _, workingHoursItem := range userWorkingHours {
			if !workingHoursItem.NotWorking &&
				workingHoursItem.StartForenoon != nil &&
				workingHoursItem.EndForenoon != nil {
				availableStart = append(availableStart, todaySetTime(start, workingHoursItem.StartForenoon))
				availableEnd = append(availableEnd, todaySetTime(start, workingHoursItem.EndForenoon))
			}

			if !workingHoursItem.NotWorking &&
				workingHoursItem.Break &&
				workingHoursItem.StartAfternoon != nil &&
				workingHoursItem.EndAfternoon != nil {
				availableStart = append(availableStart, todaySetTime(start, workingHoursItem.StartAfternoon))
				availableEnd = append(availableEnd, todaySetTime(start, workingHoursItem.EndAfternoon))
			}
		}

		appointments, exists := appointmentsByUser[user]
		if exists {
			for _, appointment := range appointments {
				// add in reverse, because appointments state unavailable dates
				availableEnd = append(availableEnd, todaySetTime(start, &appointment.Start))
				availableStart = append(availableStart, todaySetTime(start, &appointment.End))
			}
		}

		sort.Slice(availableStart, func(i, j int) bool {
			return availableStart[i].Before(availableStart[j])
		})

		sort.Slice(availableEnd, func(i, j int) bool {
			return availableEnd[i].Before(availableEnd[j])
		})

		for i := range availableStart {
			timeSlice := &TimeSlice{
				Start: availableStart[i],
				End:   availableEnd[i],
			}

			workingHoursSlices = join(
				workingHoursSlices,
				timeSlice,
			)
		}
	}

	return workingHoursSlices
}

func timeslotRanges(
	start time.Time,
	duration time.Duration,
	openingHour prisma.BranchOpeningHour,
	workingHoursSlices []*TimeSlice,
) []gqlgen.TimeslotRange {
	if openingHour.Closed == true ||
		openingHour.StartForenoon == nil ||
		openingHour.EndForenoon == nil {
		return []gqlgen.TimeslotRange{}
	}

	availableTimeSlices := calculateWorkingHoursDuringOpeningHours(start, openingHour, workingHoursSlices)

	availableStart := []time.Time{}
	availableEnd := []time.Time{}

	for _, timeSlice := range availableTimeSlices {
		availableStart = append(availableStart, timeSlice.Start)
		availableEnd = append(availableEnd, timeSlice.End)
	}

	var ranges []gqlgen.TimeslotRange

	now := time.Now()

	for i := range availableStart {
		lastAvailableSlot := availableEnd[i].Add(-duration)

		for item := availableStart[i]; item.Before(lastAvailableSlot) || item.Equal(lastAvailableSlot); item = item.Add(defaultInterval) {
			// if slot is in the past
			if item.Before(start) {
				continue
			}

			if item.Before(now) {
				continue
			}

			ranges = append(ranges, gqlgen.TimeslotRange{
				Start: item,
				End:   item.Add(duration),
			})
		}
	}

	return ranges
}
