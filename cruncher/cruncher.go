package cruncher

import (
	"github.com/pkg/errors"
	"github.com/ppxl/sagemine/core"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	dayStartTime   = "08:00:00"
	lunchStartTime = "12:00"
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
		logrus.Printf("Add new pipeline %s", redminePipeline)
		pipelineName := string(redminePipeline)

		pipeline, err := output.AddPipeline(pipelineName)
		if err != nil {
			return nil, errors.Wrap(err, "error while crunching time data")
		}

		for day, worktime := range workPerDay.WorkPerDay {
			currentDay := time.Unix(0, 0)
			if currentDay == time.Unix(0, 0) {
				firstDayString := day + "T" + dayStartTime + "Z"
				currentDay, err = time.Parse(time.RFC3339, firstDayString)
				if err != nil {
					return nil, errors.Wrapf(err, "error while creating date start time %s for pipeline %v", day, pipeline)
				}
			}

			if worktime == 0.0 {
				pipeline.PutTimeSlot(day, "-", "-")
				continue
			}

			const wallClockLayout = "15:04"
			start := currentDay.Format(wallClockLayout)
			if currentDay.Minute() != 0 {
				roundedNextHour := currentDay.Add(1 * time.Hour).Add(-time.Duration(currentDay.Minute()) * time.Minute)
				logrus.Debugf("Rounding up %s to next hour %s", start, roundedNextHour)
				start = roundedNextHour.Format(wallClockLayout)
			}
			endTime := currentDay.Add(time.Duration(worktime) * time.Hour)
			end := endTime.Format(wallClockLayout)

			pipeline.PutTimeSlot(day, start, end)
			currentDay = endTime
		}
	}

	return output, nil
}
