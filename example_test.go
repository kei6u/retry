package retry_test

import (
	"fmt"
	"time"

	"github.com/kei6u/retry"
)

func ExampleDefaultConstant() {
	retryCount := 0
	start := time.Now()
	for retry.DefaultConstant.Next() {
		retryCount++
	}
	fmt.Printf("%s elapsed after %d retries at constant intervals\n", time.Since(start), retryCount)
}

func ExampleConstant() {
	constant := retry.Constant(retry.ConstantOptions{
		Interval:    100 * time.Millisecond,
		MaxAttempts: 10,
	})
	retryCount := 0
	start := time.Now()
	for constant.Next() {
		retryCount++
	}
	fmt.Printf("%s elapsed after %d retries at constant intervals\n", time.Since(start), retryCount)
}

func ExampleDefaultExponentialBackoff() {
	retryCount := 0
	start := time.Now()
	for retry.DefaultExponentialBackoff.Next() {
		retryCount++
	}
	fmt.Printf("%s elapsed after %d retries with the exponential backoff algorithm\n", time.Since(start), retryCount)
}

func ExampleExponentialBackoff() {
	backoff := retry.ExponentialBackoff(retry.ExponentialBackoffOptions{
		BaseInterval: 100 * time.Millisecond,
		MaxInterval:  30 * time.Second,
		MaxAttempts:  10,
	})
	retryCount := 0
	start := time.Now()
	for backoff.Next() {
		retryCount++
	}
	fmt.Printf("%s elapsed after %d retries with the exponential backoff algorithm\n", time.Since(start), retryCount)
}
