package utils

import (
	"time"
)

func ElapsedTime(cbs ...func(time.Duration)) func() time.Duration {
	start := time.Now()

	return func() time.Duration {
		since := time.Since(start)

		for _, cb := range cbs {
			cb(since)
		}

		return time.Since(start)
	}
}
