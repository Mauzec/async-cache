package storage

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func emulateLoad(t *testing.T, c Cache, parallelFactor int) {
	wg := sync.WaitGroup{}

	for i := 0; i < parallelFactor; i++ {
		key := fmt.Sprintf("%d-key", i)
		value := fmt.Sprintf("%d-value", i)

		wg.Add(1)
		go func(k string, v string) {
			err := c.Set(k, v)
			assert.NoError(t, err)
			wg.Done()
		}(key, value)

		wg.Add(1)
		go func(k, v string) {
			storedValue, err := c.Get(k)
			if !errors.Is(err, ErrNotFoundByKey) {
				assert.Equal(t, v, storedValue)
			}
			wg.Done()
		}(key, value)

		wg.Add(1)
		go func(k string) {
			err := c.Delete(k)
			assert.NoError(t, err)
			wg.Done()
		}(key)
	}
	wg.Wait()
}

func emulateLoadBench(c Cache, parallelFactor int) {
	wg := sync.WaitGroup{}

	for i := 0; i < parallelFactor; i++ {
		key := fmt.Sprintf("%d-key", i)
		value := fmt.Sprintf("%d-value", i)

		wg.Add(1)
		go func(k string, v string) {
			err := c.Set(k, v)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(key, value)

		wg.Add(1)
		go func(k, v string) {
			_, err := c.Get(k)
			if err != nil && !errors.Is(err, ErrNotFoundByKey) {
				panic(err)
			}
			wg.Done()
		}(key, value)

		wg.Add(1)
		go func(k string) {
			err := c.Delete(k)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(key)
	}
	wg.Wait()
}

func emulateLoadWithMetrics(t *testing.T, cm CacheWithMetrics, parallelFactor int) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		emulateLoad(t, cm, parallelFactor)
		wg.Done()
	}()

	var min, max int64
	for i := 0; i < parallelFactor; i++ {
		wg.Add(1)
		go func() {
			total := cm.TotalAmount()
			if total > max {
				max = total
			}
			if total < min {
				min = total
			}
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(min, max)
}
