package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const theDate = "2021-05-05"

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
