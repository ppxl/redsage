package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const theDate = "2021-05-05"
const startTime = "08:00"
const endTime = "09:00"

func TestRedmineWorkPerDay_PutWorkTime(t *testing.T) {
	t.Run("should add another work time", func(t *testing.T) {
		sut := newRedmineWorkPerDay()
		sut.PutWorkTime(theDate, 4.75)
		assert.Equal(t, 4.75, sut.WorkTime(theDate))
		sut.PutWorkTime(theDate, 1.25)
		assert.Equal(t, 6.0, sut.WorkTime(theDate))
	})
}

func TestRedmineWorkPerDay_WorkTime(t *testing.T) {
	t.Run("should return 0 for new work time", func(t *testing.T) {
		sut := newRedmineWorkPerDay()
		assert.Equal(t, 0.0, sut.WorkTime(theDate))
	})
	t.Run("should return given value for new work time", func(t *testing.T) {
		sut := newRedmineWorkPerDay()
		sut.PutWorkTime(theDate, 4.75)
		assert.Equal(t, 4.75, sut.WorkTime(theDate))
	})
}

func TestSageWorkPerDay_PutTimeSlot(t *testing.T) {
	t.Run("should create new time slot", func(t *testing.T) {
		sut := SageWorkPerDay{}

		// when
		sut.PutTimeSlot(theDate, startTime, endTime)
		actual := sut.timeSlot(theDate)

		// then
		require.NotNil(t, actual)
		expected := &TimeSlot{
			Start: startTime,
			End:   endTime,
		}
		assert.Equal(t, expected, actual)
	})
	t.Run("should add another, non-conflicting time slot to the same date", func(t *testing.T) {
		sut := SageWorkPerDay{}

		// when
		sut.PutTimeSlot(theDate, startTime, endTime)
		actual := sut.timeSlot(theDate)

		// then
		require.NotNil(t, actual)
		expected := &TimeSlot{
			Start: startTime,
			End:   endTime,
		}
		assert.Equal(t, expected, actual)
	})
}

func TestSageWorkPerDay_timeSlot(t *testing.T) {
	t.Run("should return nil for unset time slot", func(t *testing.T) {
		sut := SageWorkPerDay{}

		actual := sut.timeSlot(theDate)

		assert.Nil(t, actual)
	})
	t.Run("should return set time slot", func(t *testing.T) {
		sut := SageWorkPerDay{}

		sut.PutTimeSlot(theDate, startTime, endTime)

		// when
		actual := sut.timeSlot(theDate)

		// then
		require.NotNil(t, actual)
		expected := &TimeSlot{
			Start: startTime,
			End:   endTime,
		}
		assert.Equal(t, expected, actual)
	})
}
