package cruncher

import (
	"github.com/ppxl/sagemine/core"
)

const (
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

	return nil, nil
}
