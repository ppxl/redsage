package cruncher

import (
	"github.com/ppxl/sagemine/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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
		expectedPipelineA.PutTimeSlot(date4, "08:00", "09:00")
		expectedPipelineA.PutTimeSlot(date5, "08:00", "12:00")
		expectedPipelineA.PutTimeSlot(date6, "08:00", "12:00")
		expectedPipelineA.PutTimeSlot(date6, "13:00", "13:30")
		expectedPipelineA.PutTimeSlot(date7, "08:00", "12:00")
		expectedPipelineA.PutTimeSlot(date7, "13:00", "15:00")
		assert.Equal(t, 5, actual.NamedDaySageValues[pipelineAName].Days())
		assert.Equal(t, (*expectedPipelineA)[date3], (*actual.NamedDaySageValues[pipelineAName])[date3])
		assert.Equal(t, (*expectedPipelineA)[date4], (*actual.NamedDaySageValues[pipelineAName])[date4])
		assert.Equal(t, (*expectedPipelineA)[date5], (*actual.NamedDaySageValues[pipelineAName])[date5])
		assert.Equal(t, (*expectedPipelineA)[date6], (*actual.NamedDaySageValues[pipelineAName])[date6])
		assert.Equal(t, (*expectedPipelineA)[date7], (*actual.NamedDaySageValues[pipelineAName])[date7])
	})
	t.Run("should add 45 minutes lunch break to joined pipeline", func(t *testing.T) {
		input := core.NewPipelineData()
		pipelineA, _ := input.AddPipeline(pipelineAName)
		pipelineA.PutWorkTime(date3, 0)
		pipelineA.PutWorkTime(date4, 1)
		pipelineA.PutWorkTime(date5, 4)
		pipelineA.PutWorkTime(date6, 4.5)
		pipelineA.PutWorkTime(date7, 6)

		config := Config{
			LunchBreakInMin: 45,
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
		expectedPipelineA.PutTimeSlot(date4, "08:00", "09:00")
		expectedPipelineA.PutTimeSlot(date5, "08:00", "12:00")
		expectedPipelineA.PutTimeSlot(date6, "08:00", "12:00")
		expectedPipelineA.PutTimeSlot(date6, "12:45", "13:15")
		expectedPipelineA.PutTimeSlot(date7, "08:00", "12:00")
		expectedPipelineA.PutTimeSlot(date7, "12:45", "14:45")
		assert.Equal(t, 5, actual.NamedDaySageValues[pipelineAName].Days())
		assert.Equal(t, (*expectedPipelineA)[date3], (*actual.NamedDaySageValues[pipelineAName])[date3])
		assert.Equal(t, (*expectedPipelineA)[date4], (*actual.NamedDaySageValues[pipelineAName])[date4])
		assert.Equal(t, (*expectedPipelineA)[date5], (*actual.NamedDaySageValues[pipelineAName])[date5])
		assert.Equal(t, (*expectedPipelineA)[date6], (*actual.NamedDaySageValues[pipelineAName])[date6])
		assert.Equal(t, (*expectedPipelineA)[date7], (*actual.NamedDaySageValues[pipelineAName])[date7])
	})
}

func Test_endTimeFallsIntoLunch(t *testing.T) {
	t.Run("should be false for 11:00 < 12:00", func(t *testing.T) {
		worktimeEnd, _ := time.Parse(time.RFC3339, date5+"T11:00:00Z")
		actualDiff, actualHit := endTimeFallsIntoLunch(worktimeEnd, date5)
		assert.False(t, actualHit)
		assert.Equal(t, time.Duration(0), actualDiff)
	})
	t.Run("should be false for 12:00 == 12:00", func(t *testing.T) {
		worktimeEnd, _ := time.Parse(time.RFC3339, date5+"T12:00:00Z")
		actualDiff, actualHit := endTimeFallsIntoLunch(worktimeEnd, date5)
		assert.False(t, actualHit)
		assert.Equal(t, time.Duration(0), actualDiff)
	})
	t.Run("should be true for 12:01 > 12:00", func(t *testing.T) {
		worktimeEnd, _ := time.Parse(time.RFC3339, date5+"T12:01:00Z")
		actualDiff, actualMiss := endTimeFallsIntoLunch(worktimeEnd, date5)
		assert.True(t, actualMiss)
		assert.Equal(t, 1*time.Minute, actualDiff)
	})
	t.Run("should be true for 13:00 > 12:00", func(t *testing.T) {
		worktimeEnd, _ := time.Parse(time.RFC3339, date5+"T13:00:00Z")
		actualDiff, actualMiss := endTimeFallsIntoLunch(worktimeEnd, date5)
		assert.True(t, actualMiss)
		assert.Equal(t, 1*time.Hour, actualDiff)
	})
}
