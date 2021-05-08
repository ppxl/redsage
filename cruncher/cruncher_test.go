package cruncher

import (
	"github.com/ppxl/sagemine/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_cruncher_Crunch(t *testing.T) {
	t.Run("should join pipelines per day", func(t *testing.T) {
		input := core.NewPipelineData()
		pipelineA := core.RedmineWorkPerDay{}
		pipelineA["2021-05-03"] = 1.50
		pipelineA["2021-05-04"] = 4
		pipelineA["2021-05-05"] = 0
		pipelineA["2021-05-06"] = 4.15
		pipelineA["2021-05-07"] = 8.25

		pipelineB := core.RedmineWorkPerDay{}
		pipelineB["2021-05-03"] = 0
		pipelineB["2021-05-04"] = 4.25
		pipelineB["2021-05-05"] = 0
		pipelineB["2021-05-06"] = 3.75
		pipelineB["2021-05-06"] = 0

		input.NamedDayRedmineValues["Pipeline A"] = &pipelineA

		config := Config{
			LunchBreakInMin:     60,
			SinglePipelineNames: nil,
		}

		sut := New()

		//when
		actual, err := sut.Crunch(input, config)

		// then
		require.NoError(t, err)
		expected := core.CrunchedOutput{}
		assert.Equal(t, expected, actual)
	})
}
