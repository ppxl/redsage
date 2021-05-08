package cruncher

import (
	"github.com/ppxl/sagemine/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	pipelineAName = "Pipeline A"
	date3         = "2021-05-03"
	date4         = "2021-05-04"
	date5         = "2021-05-05"
	date6         = "2021-05-06"
	date7         = "2021-05-07"
)

func Test_cruncher_Crunch(t *testing.T) {
	t.Run("should add 1 hour lunch break to joined pipeline", func(t *testing.T) {
		input := core.NewPipelineData()
		pipelineA, _ := input.AddPipeline(pipelineAName)
		pipelineA.PutWorkTime(date3, 0)
		pipelineA.PutWorkTime(date4, 1)
		pipelineA.PutWorkTime(date5, 4)
		pipelineA.PutWorkTime(date6, 4.5)
		pipelineA.PutWorkTime(date7, 6)

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
		expectedPipelineA.PutTimeSlot(date3, "-", "-")

		expectedPipelineA.PutTimeSlot(date4, "08:00", "12:00")
		expectedPipelineA.PutTimeSlot(date4, "12:00", "14:00")

		expectedPipelineA.PutTimeSlot(date5, "08:00", "12:00")
		expectedPipelineA.PutTimeSlot(date5, "13:00", "17:00")

		expectedPipelineA.PutTimeSlot(date6, "08:00", "12:00")
		expectedPipelineA.PutTimeSlot(date6, "13:00", "16:30")

		expectedPipelineA.PutTimeSlot(date7, "08:00", "12:00")
		expectedPipelineA.PutTimeSlot(date7, "13:00", "15:00")
		assert.Equal(t, 5, actual.NamedDaySageValues[pipelineAName].Days())
		assert.ElementsMatch(t, (*expectedPipelineA)[date3], (*actual.NamedDaySageValues[pipelineAName])[date3])
		assert.ElementsMatch(t, (*expectedPipelineA)[date4], (*actual.NamedDaySageValues[pipelineAName])[date4])
		assert.ElementsMatch(t, (*expectedPipelineA)[date5], (*actual.NamedDaySageValues[pipelineAName])[date5])
		assert.ElementsMatch(t, (*expectedPipelineA)[date6], (*actual.NamedDaySageValues[pipelineAName])[date6])
		assert.ElementsMatch(t, (*expectedPipelineA)[date7], (*actual.NamedDaySageValues[pipelineAName])[date7])
	})
}
