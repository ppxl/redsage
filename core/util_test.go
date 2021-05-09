package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestParseDateWithTime(t *testing.T) {
	t.Run("should parse time", func(t *testing.T) {
		actual, err := ParseDateWithTime("2021-05-05", "17:00:00")

		require.NoError(t, err)
		expected := time.Date(2021, 05, 05, 17, 0, 0, 0, time.UTC)
		assert.Equal(t, expected, actual)
	})
	t.Run("should error for missing seconds", func(t *testing.T) {
		_, err := ParseDateWithTime("2021-05-05", "17:00")

		require.Error(t, err)
	})
}
