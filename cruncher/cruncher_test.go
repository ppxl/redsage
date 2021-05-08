package cruncher

import (
	"github.com/ppxl/sagemine/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	pipelineAName = "Pipeline A"
	pipelineBName = "Pipeline 2/B"
)

func Test_cruncher_Crunch(t *testing.T) {
	t.Run("should add 1 hour lunch break to joined pipeline", func(t *testing.T) {
		input := core.NewPipelineData()
		pipelineA, _ := input.AddPipeline(pipelineAName)
		pipelineA.PutWorkTime("2021-05-03", 0)
		pipelineA.PutWorkTime("2021-05-04", 1)
		pipelineA.PutWorkTime("2021-05-05", 4)
		pipelineA.PutWorkTime("2021-05-06", 4.5)
		pipelineA.PutWorkTime("2021-05-07", 6)

		config := Config{
			LunchBreakInMin: 60,
		}

		sut := New()

		//when
		actual, err := sut.Crunch(input, config)

		// then
		require.NoError(t, err)
		expected := core.NewCrunchedOutput()
		expectedPipelineA, err := expected.AddPipeline(pipelineAName)
		require.NoError(t, err)
		expectedPipelineA.PutTimeSlot("2021-05-03", "-", "-")

		expectedPipelineA.PutTimeSlot("2021-05-04", "08:00", "12:00")
		expectedPipelineA.PutTimeSlot("2021-05-04", "13:00", "14:00")

		expectedPipelineA.PutTimeSlot("2021-05-05", "08:00", "12:00")
		expectedPipelineA.PutTimeSlot("2021-05-05", "13:00", "17:00")

		expectedPipelineA.PutTimeSlot("2021-05-06", "08:00", "12:00")
		expectedPipelineA.PutTimeSlot("2021-05-06", "13:00", "16:30")

		expectedPipelineA.PutTimeSlot("2021-05-07", "08:00", "12:00")
		expectedPipelineA.PutTimeSlot("2021-05-07", "13:00", "15:00")
		assert.Equal(t, 4, actual.NamedDaySageValues[pipelineAName].Days())
		assert.Equal(t, expected.String(), actual.String())

	})
}
