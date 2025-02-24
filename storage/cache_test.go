package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Cache(t *testing.T) {
	t.Parallel()

	cache := NewCacheWithMetrics()

	t.Run("correctly stored value", func(t *testing.T) {
		t.Parallel()
		key := "1"
		value := "one"

		err := cache.Set(key, value)
		assert.NoError(t, err)

		storedValue, err := cache.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, value, storedValue)

	})

	t.Run("correctly update value", func(t *testing.T) {
		t.Parallel()
		key := "1"
		value := "one"

		err := cache.Set(key, value)
		assert.NoError(t, err)

		storedValue, err := cache.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, value, storedValue)

		err = cache.Set("2", "two")
		assert.NoError(t, err)

		storedValue, err = cache.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, value, storedValue)

	})

	t.Run("no data races", func(t *testing.T) {
		t.Parallel()

		parallelFactor := 100_000
		emulateLoad(t, cache, parallelFactor)
	})

}
