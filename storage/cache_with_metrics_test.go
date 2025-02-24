package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CacheWithMetrics(t *testing.T) {
	t.Parallel()

	testCache := NewCacheWithMetrics()

	t.Run("correctly stored value", func(t *testing.T) {
		t.Parallel()

		key := "1"
		value := "one"

		err := testCache.Set(key, value)
		assert.NoError(t, err)
		storedValue, err := testCache.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, value, storedValue)
	})

	t.Run("no data races", func(t *testing.T) {
		t.Parallel()

		parallelFactor := 100_000
		emulateLoadWithMetrics(t, testCache, parallelFactor)
	})
}
