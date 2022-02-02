package retry_test

import (
	"fmt"
	"time"

	"github.com/kei6u/retry"
)

func ExampleConstant() {
	constant := retry.Constant(retry.ConstantOptions{
		Interval:    100 * time.Millisecond,
		MaxAttempts: 10,
	})
	retryCount := 0
	start := time.Now()
	for constant.Next() {
		retryCount++
		fmt.Printf("retry %d: %s elapsed\n", retryCount, time.Since(start))
		start = time.Now()
	}
}

func ExampleExponentialBackoff() {
	backoff := retry.ExponentialBackoff(retry.ExponentialBackoffOptions{
		BaseInterval: 100 * time.Millisecond,
		MaxInterval:  10 * time.Second,
		MaxAttempts:  10,
	})
	retryCount := 0
	start := time.Now()
	for backoff.Next() {
		retryCount++
		fmt.Printf("retry %d: %s elapsed\n", retryCount, time.Since(start))
		start = time.Now()
	}
}
