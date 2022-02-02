package retry

import (
	"context"
	"testing"
	"time"
)

func TestDefaultConstantRetry(t *testing.T) {
	t.Parallel()
	defaultContextWithTimeout = func() (context.Context, context.CancelFunc) {
		return context.WithTimeout(context.TODO(), 10*time.Millisecond)
	}
	start := time.Now()
	for DefaultConstant.Next() {
	}
	if time.Since(start) < 10*time.Millisecond {
		t.Fatalf("expected to timeout after 100ms")
	}
}

func TestConstantRetryWithMaxAttempts(t *testing.T) {
	t.Parallel()
	retryCount := 0
	constant := Constant(ConstantOptions{
		Interval:    time.Millisecond,
		MaxAttempts: 10,
	})
	for constant.Next() {
		retryCount++
	}
	if retryCount != 10 {
		t.Fatalf("expected 10 retries, but %d", retryCount)
	}
}

func TestConstantRetryWithTimeoutContext(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()
	constant := Constant(ConstantOptions{
		Context:  ctx,
		Interval: time.Second,
	})
	start := time.Now()
	for constant.Next() {
	}
	if time.Since(start) < 10*time.Millisecond {
		t.Fatal("expected timeout after 10ms")
	}
}

func TestConstantRetryMaxAttemptsAndWithTimeoutContext(t *testing.T) {
	t.Parallel()
	t.Run("prioritize max attempts", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		defer cancel()
		constant := Constant(ConstantOptions{
			Context:     ctx,
			Interval:    time.Millisecond,
			MaxAttempts: 10,
		})
		start := time.Now()
		retryCount := 0
		for constant.Next() {
			retryCount++
		}
		if time.Second < time.Since(start) && retryCount < 10 {
			t.Fatal("expected to reach 10 retries before 1 second elapsed")
		}
	})
	t.Run("prioritize timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
		defer cancel()
		constant := Constant(ConstantOptions{
			Context:     ctx,
			Interval:    2 * time.Millisecond,
			MaxAttempts: 10,
		})
		start := time.Now()
		retryCount := 0
		for constant.Next() {
			retryCount++
		}
		if 10 <= retryCount && 10*time.Millisecond < time.Since(start) {
			t.Logf("%d retries took %s", retryCount, time.Since(start))
			t.Fatal("expected to timeout before 10 retries")
		}
	})
}

func TestExponentialBackoffRetry(t *testing.T) {
	t.Parallel()
	defaultContextWithTimeout = func() (context.Context, context.CancelFunc) {
		return context.WithTimeout(context.TODO(), 10*time.Millisecond)
	}
	start := time.Now()
	for DefaultExponentialBackoff.Next() {
	}
	if time.Since(start) < 10*time.Millisecond {
		t.Fatalf("expected to timeout after 10ms")
	}
}

func TestExponentialBackoffRetryWithMaxAttempts(t *testing.T) {
	t.Parallel()
	backoff := ExponentialBackoff(ExponentialBackoffOptions{
		BaseInterval: time.Millisecond,
		MaxAttempts:  10,
	})
	start := time.Now()
	retryCount := 0
	for backoff.Next() {
		retryCount++
	}
	if retryCount != 10 && time.Since(start) < 600*time.Millisecond {
		t.Fatal("expected 10 retries and elapse at least 600ms")
	}
}

func TestExponentialBackoffRetryWithTimeout(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
	defer cancel()
	backoff := ExponentialBackoff(ExponentialBackoffOptions{
		Context:      ctx,
		BaseInterval: time.Second,
	})
	start := time.Now()
	for backoff.Next() {
	}
	if time.Since(start) < 10*time.Millisecond {
		t.Fatal("expected elapse at least 10ms")
	}
}

func TestExponentialBackoffRetryMaxAttemptsAndWithTimeoutContext(t *testing.T) {
	t.Parallel()
	t.Run("prioritize max attempts", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
		defer cancel()
		backoff := ExponentialBackoff(ExponentialBackoffOptions{
			Context:      ctx,
			BaseInterval: time.Millisecond,
			MaxAttempts:  10,
		})
		start := time.Now()
		retryCount := 0
		for backoff.Next() {
			retryCount++
		}
		if time.Second < time.Since(start) && retryCount < 10 {
			t.Fatal("expected to reach 10 retries before 1 second elapsed")
		}
	})
	t.Run("prioritize timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Millisecond)
		defer cancel()
		backoff := ExponentialBackoff(ExponentialBackoffOptions{
			Context:      ctx,
			BaseInterval: time.Millisecond,
			MaxAttempts:  10,
		})
		start := time.Now()
		retryCount := 0
		for backoff.Next() {
			retryCount++
		}
		if 10 <= retryCount && 10*time.Millisecond < time.Since(start) {
			t.Fatal("expected to timeout before 10 retries")
		}
	})
}
