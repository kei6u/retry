# retry

[![.github/workflows/test.yaml](https://github.com/kei6u/retry/actions/workflows/test.yaml/badge.svg)](https://github.com/kei6u/retry/actions/workflows/test.yaml)
[![GoDoc](https://godoc.org/github.com/kei6u/retry?status.svg&style=flat-square)](http://godoc.org/github.com/kei6u/retry)

This Go library is made from only standard libraries and provides retry functionality for general operations.
You can choose a retry algorithm from constant intervals, decorrelated jitter algorithm, exponential backoff algorithm.

## Motivation

I would like to use retry algorithms inspired by [Exponential Backoff And Jitter | AWS Architecture Blog](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/) in Go.

## When

- When the database returns internal errors, your application will retry to run SQL with this library.
- When the API returns 5xx errors, your application will retry to call it with this library.

## Usage

See the [document](https://pkg.go.dev/github.com/kei6u/retry) and run [examples](https://pkg.go.dev/github.com/kei6u/retry#pkg-examples).

```bash
go get github.com/kei6u/retry
```

```go
import "github.com/kei6u/retry"
```

```go
r := retry.New(retry.Jitter{})
for r.Next() {
	resp, err := http.Get("http://example.com")
	if err != nil {
		_ = resp.Body.Close()
		return err
	}
	if 500 <= resp.StatusCode && resp.StatusCode < 600 {
		_ = resp.Body.Close()
		continue
	}
	if resp.StatusCode == 200 {
		_ = resp.Body.Close()
		return nil
	}
}
```

## Algorithms

### Jitter (Recommended)

This algorithm provides retries with "Decorrelated Jitter" from [Exponential Backoff And Jitter | AWS Architecture Blog](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/). This blog introduces this algorithm as better. You can run the [example](https://pkg.go.dev/github.com/kei6u/retry#example-Jitter) on your browser.

> It’s worth noting that none of these approaches fundamentally change the N2 nature of the work to be done, but do substantially reduce work at reasonable levels of contention. The return on implementation complexity of using jittered backoff is huge, and it should be considered a standard approach for remote clients.
> https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/

### Constant

This algorithm provides retries at constant intervals. You can run the [example](https://pkg.go.dev/github.com/kei6u/retry#example-Constant) on your browser.

### Exponential backoff

This algorithm provides retries with the exponential backoff algorithm. You can run the [example](https://pkg.go.dev/github.com/kei6u/retry#example-ExponentialBackoff) on your browser.

> Exponential backoff is an algorithm that uses feedback to multiplicatively decrease the rate of some process, in order to gradually find an acceptable rate. These algorithms find usage in a wide range of systems and processes, with radio networks and computer networks being particularly notable.
> https://en.wikipedia.org/wiki/Exponential_backoff
