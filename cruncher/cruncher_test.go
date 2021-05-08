package cruncher

import (
	"github.com/ppxl/sagemine/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	pipelineAName       = "Pipeline A"
	pipelineAJoinedName = "Pipeline A-joined"
	pipelineBName       = "Pipeline 2/B"
)

func Test_cruncher_Crunch(t *testing.T) {
	t.Run("should join two pipelines and add their values per day", func(t *testing.T) {
		input := core.NewPipelineData()
		pipelineA, _ := input.AddPipeline(pipelineAName)
		pipelineA.PutWorkTime("2021-05-03", 0)
		pipelineA.PutWorkTime("2021-05-04", 0)
		pipelineA.PutWorkTime("2021-05-05", 1)
		pipelineA.PutWorkTime("2021-05-06", 2)
		pipelineA.PutWorkTime("2021-05-07", 3.5)

		pipelineB, _ := input.AddPipeline(pipelineBName)
		pipelineB.PutWorkTime("2021-05-03", 0)
		pipelineB.PutWorkTime("2021-05-04", 1)
		pipelineB.PutWorkTime("2021-05-05", 2)
		pipelineB.PutWorkTime("2021-05-06", 1.5)
		pipelineB.PutWorkTime("2021-05-07", 0.5)

		config := Config{
			LunchBreakInMin:     60,
			SinglePipelineNames: nil,
		}

		sut := New()

		//when
		actual, err := sut.Crunch(input, config)

		// then
		require.NoError(t, err)
		expected := core.NewCrunchedOutput()
		expectedPipelineA, err := expected.AddPipeline(pipelineAJoinedName)
		require.NoError(t, err)
		// no timeslot for 2021-05-03
		expectedPipelineA.PutTimeSlot("2021-05-04", "08:00", "09:00")
		expectedPipelineA.PutTimeSlot("2021-05-05", "08:00", "11:00")
		expectedPipelineA.PutTimeSlot("2021-05-06", "08:00", "11:30")
		expectedPipelineA.PutTimeSlot("2021-05-07", "08:00", "12:00")
		assert.Equal(t, expected, actual)
	})
	t.Run("should add 1 hour lunch break to joined pipeline", func(t *testing.T) {
		input := core.NewPipelineData()
		pipelineA, _ := input.AddPipeline(pipelineAName)
		pipelineA.PutWorkTime("2021-05-03", 0)
		pipelineA.PutWorkTime("2021-05-04", 1)
		pipelineA.PutWorkTime("2021-05-05", 4)
		pipelineA.PutWorkTime("2021-05-06", 6)
		pipelineA.PutWorkTime("2021-05-07", 4.5)

		pipelineB, _ := input.AddPipeline(pipelineBName)
		pipelineB.PutWorkTime("2021-05-03", 9)
		pipelineB.PutWorkTime("2021-05-04", 4)
		pipelineB.PutWorkTime("2021-05-05", 4)
		pipelineB.PutWorkTime("2021-05-06", 1.5)
		pipelineB.PutWorkTime("2021-05-07", 0.5)

		config := Config{
			LunchBreakInMin:     60,
			SinglePipelineNames: nil,
		}

		sut := New()

		//when
		actual, err := sut.Crunch(input, config)

		// then
		require.NoError(t, err)
		expected := core.NewCrunchedOutput()
		expectedPipelineA, err := expected.AddPipeline(pipelineAJoinedName)
		require.NoError(t, err)
		expectedPipelineA.PutTimeSlot("2021-05-03", "08:00", "12:00")
		expectedPipelineA.PutTimeSlot("2021-05-03", "13:00", "18:00")

		expectedPipelineA.PutTimeSlot("2021-05-04", "08:00", "12:00")
		expectedPipelineA.PutTimeSlot("2021-05-04", "13:00", "14:00")

		expectedPipelineA.PutTimeSlot("2021-05-05", "08:00", "12:00")
		expectedPipelineA.PutTimeSlot("2021-05-05", "13:00", "17:00")

		expectedPipelineA.PutTimeSlot("2021-05-06", "08:00", "12:00")
		expectedPipelineA.PutTimeSlot("2021-05-06", "13:00", "16:30")

		expectedPipelineA.PutTimeSlot("2021-05-07", "08:00", "12:00")
		expectedPipelineA.PutTimeSlot("2021-05-07", "13:00", "15:00")
		assert.Equal(t, expected, actual)
	})
}
