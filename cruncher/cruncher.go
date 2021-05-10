package cruncher

import (
	"github.com/pkg/errors"
	"github.com/ppxl/sagemine/core"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	dayStartTime    = "08:00:00"
	lunchStartTime  = "12:00:00"
	wallClockLayout = "15:04"
)

// Config contains configuration values that modify the number crunching behaviour.
type Config struct {
	LunchBreakInMin int
}

// Cruncher provides methods for transforming values from a redmine pipeline data.
type Cruncher interface {
	// Crunch executes the merging and splitting values from a redmine pipeline data and returns them.
	Crunch(pdata *core.PipelineData, config Config) (*core.CrunchedOutput, error)
}

type cruncher struct {
}

func New() *cruncher {
	return &cruncher{}
}

// Crunch executes the merging and splitting values from a CSV file and prints the output in Sage-relatable manner.
func (c *cruncher) Crunch(pdata *core.PipelineData, config Config) (*core.CrunchedOutput, error) {
	output := core.NewCrunchedOutput()

	for redminePipeline, workPerDay := range pdata.NamedDayRedmineValues {
		logrus.Debugf("Add new pipeline %s", redminePipeline)
		pipelineName := string(redminePipeline)

		pipeline, err := output.AddPipeline(pipelineName)
		if err != nil {
			return nil, errors.Wrap(err, "error while crunching time data")
		}

		for day, worktime := range workPerDay.WorkPerDay {
			firstDayString := day + "T" + dayStartTime + "Z"
			currentDayAndTime, err := time.Parse(time.RFC3339, firstDayString)
			if err != nil {
				return nil, errors.Wrapf(err, "error while creating date start time %s for pipeline %v", day, pipeline)
			}

			if containsNoWorkTime(worktime) {
				pipeline.PutEmptyTimeSlot(day)
				continue
			}

			start := currentDayAndTime.Format(wallClockLayout)
			if !startsWorkAtOClock(currentDayAndTime) {
				roundedNextHour := roundWorkTimeToNextHour(currentDayAndTime)
				logrus.Debugf("Rounding up %s to next hour: %s", start, roundedNextHour)
				currentDayAndTime = roundedNextHour
				start = currentDayAndTime.Format(wallClockLayout)
			}

			// decimals don't work well with duration: Do instead manual minute calculation
			calcedEndTime := time.Duration(worktime*60) * time.Minute
			endTime := currentDayAndTime.Add(calcedEndTime)

			diff, endTimeIntersectsWithLunchtime := endTimeFallsIntoLunch(endTime, day)
			if endTimeIntersectsWithLunchtime {
				logrus.Debugf("Time slot %s - %s falls into lunch by %s. Breaking up into two parts...", currentDayAndTime, endTime, diff)
				endBeforeLunch := endTime.Add(-diff)
				pipeline.PutTimeSlot(day, start, endBeforeLunch.Format(wallClockLayout))
				// update current time to enable correct timing of the second slot
				currentDayAndTime = endBeforeLunch

				startAfterLunch := currentDayAndTime.Add(time.Duration(config.LunchBreakInMin) * time.Minute)
				endAfterLunch := startAfterLunch.Add(diff)
				pipeline.PutTimeSlot(day, startAfterLunch.Format(wallClockLayout), endAfterLunch.Format(wallClockLayout))
				currentDayAndTime = endAfterLunch
			} else {
				end := endTime.Format(wallClockLayout)
				pipeline.PutTimeSlot(day, start, end)
				currentDayAndTime = endTime
			}

		}
	}

	return output, nil
}

// endTimeFallsIntoLunch returns false if the given time does not overlap with to configured lunch time. Otherwise true
// and the (aboslute) duration of the overlap will be returned.
func endTimeFallsIntoLunch(workTimeEnd time.Time, day string) (time.Duration, bool) {
	parsedLunchtime, err := time.Parse(time.RFC3339, day+"T"+lunchStartTime+"Z")
	if err != nil {
		panic("lunchtime: " + err.Error())
	}

	fallsIntoLunchTime := workTimeEnd.After(parsedLunchtime)
	if fallsIntoLunchTime {
		return workTimeEnd.Sub(parsedLunchtime), true
	}

	return 0, false
}

func roundWorkTimeToNextHour(currentWorkTime time.Time) time.Time {
	// go one hour forward and subtract the actual minutes to get the full next hour
	roundedNextHour := currentWorkTime.Add(1 * time.Hour).Add(-time.Duration(currentWorkTime.Minute()) * time.Minute)
	return roundedNextHour
}

func startsWorkAtOClock(workTime time.Time) bool {
	return workTime.Minute() == 0
}

func containsNoWorkTime(worktime float64) bool {
	return worktime == 0.0
}
