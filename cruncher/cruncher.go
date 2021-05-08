package cruncher

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/ppxl/sagemine/core"
	"time"
)

const (
	dayStartTime   = "08:00:00"
	lunchStartTime = "12:00"
)

// Config contains configuration values that modify the number crunching behaviour.
type Config struct {
	LunchBreakInMin     int
	SinglePipelineNames []core.PipelineName
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

	var joinedPipeline *core.SageWorkPerDay
	var err error

	for redminePipeline, workPerDay := range pdata.NamedDayRedmineValues {
		fmt.Printf("%s, \t", redminePipeline)

		var pipeline *core.SageWorkPerDay
		if true {
			if joinedPipeline == nil {
				joinedPipeline, err = output.AddPipeline(string(redminePipeline) + "-joined")

				if err != nil {
					return nil, errors.Wrap(err, "error while crunching time data")
				}
			}
			pipeline = joinedPipeline
		} else {
			pipeline, err = output.AddPipeline(string(redminePipeline))

			if err != nil {
				return nil, errors.Wrap(err, "error while crunching time data")
			}
		}

		currentDay := time.Unix(0, 0)
		for day, worktime := range (map[string]float64)(*workPerDay) {
			if currentDay == time.Unix(0, 0) {
				firstDayString := day + "T" + dayStartTime + "Z"
				currentDay, err = time.Parse(time.RFC3339, firstDayString)
				if err != nil {
					return nil, errors.Wrapf(err, "error while creating date start time %s for pipeline %v", day, pipeline)
				}
			}

			fmt.Printf("%s, %0.2f\t", redminePipeline, worktime)

			start := currentDay.Format("15:04")
			currentDay.Add(time.Duration(worktime))
			end := currentDay.Format("15:04")

			pipeline.PutTimeSlot(day, start, end)
		}
		fmt.Println()
	}
	return output, nil
}
