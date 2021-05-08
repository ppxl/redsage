package cruncher

import (
	"fmt"
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
	var joinedPipelineName string
	var err error

	for redminePipeline, workPerDay := range pdata.NamedDayRedmineValues {
		var pipeline *core.SageWorkPerDay
		if true {
			if joinedPipeline == nil {
				pipelineName := string(redminePipeline) + "-joined"
				joinedPipeline, err = output.AddPipeline(pipelineName)
				joinedPipelineName = pipelineName

				logrus.Printf("Add new pipeline %s", redminePipeline)

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

		fmt.Printf("%s, \t", joinedPipelineName)
		currentDay := time.Unix(0, 0)

		for day, worktime := range (map[string]float64)(*workPerDay) {
			if currentDay == time.Unix(0, 0) {
				firstDayString := day + "T" + dayStartTime + "Z"
				currentDay, err = time.Parse(time.RFC3339, firstDayString)
				if err != nil {
					return nil, errors.Wrapf(err, "error while creating date start time %s for pipeline %v", day, pipeline)
				}
			}

			fmt.Printf("%s, %0.2f", day, worktime)
			if worktime == 0.0 {
				fmt.Print("*\t")
				continue
			}

			fmt.Print("\t")

			start := currentDay.Format("15:04")
			currentDay.Add(time.Duration(worktime) * time.Hour)
			end := currentDay.Format("15:04")

			pipeline.PutTimeSlot(day, start, end)
		}

		fmt.Println()
	}
	return output, nil
}
