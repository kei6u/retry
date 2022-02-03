package retry_test

import (
	"fmt"
	"time"

	"github.com/kei6u/retry"
)

func ExampleConstant() {
	retrier := retry.New(retry.Constant{
		Interval:    100 * time.Millisecond,
		MaxAttempts: 10,
	})
	retryCount := 0
	start := time.Now()
	for retrier.Next() {
		retryCount++
		fmt.Printf("retry %d: %s elapsed\n", retryCount, time.Since(start))
		start = time.Now()
	}
}

func ExampleExponentialBackoff() {
	retrier := retry.New(retry.ExponentialBackoff{
		BaseInterval: 100 * time.Millisecond,
		MaxInterval:  10 * time.Second,
		MaxAttempts:  10,
	})
	retryCount := 0
	start := time.Now()
	for retrier.Next() {
		retryCount++
		fmt.Printf("retry %d: %s elapsed\n", retryCount, time.Since(start))
		start = time.Now()
	}
}
