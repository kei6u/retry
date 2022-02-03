# retry

[![.github/workflows/test.yaml](https://github.com/kei6u/retry/actions/workflows/test.yaml/badge.svg)](https://github.com/kei6u/retry/actions/workflows/test.yaml)
[![GoDoc](https://godoc.org/github.com/kei6u/retry?status.svg&style=flat-square)](http://godoc.org/github.com/kei6u/retry)

This Go library is made from only standard libraries and provides retry functionality for general operations.
You can choose a retry strategy from constant intervals or the exponential backoff algorithm.

## Motivation

There are popular fantastic similar libraries already.
However, I want a new library that provides a more straightforward interface and implementation. This is the biggest and only motivation to create this library.

## Usage

See the [document](https://pkg.go.dev/github.com/kei6u/retry) and run [examples](https://pkg.go.dev/github.com/kei6u/retry#pkg-examples).

### Import

```bash
go get github.com/kei6u/retry
```

```go
import "github.com/kei6u/retry"
```

### Constant

This strategy provides retries at constant intervals. You can run the [example](https://pkg.go.dev/github.com/kei6u/retry#example-Constant) on your browser.

### Exponential backoff

This strategy provides retries with the exponential backoff algorithm. You can run the [example](https://pkg.go.dev/github.com/kei6u/retry#example-ExponentialBackoff) on your browser.

"Exponential backoff is an algorithm that uses feedback to multiplicatively decrease the rate of some process, in order to gradually find an acceptable rate. These algorithms find usage in a wide range of systems and processes, with radio networks and computer networks being particularly notable." ("Exponential backoff," n.d.)

Reference

Exponential backoff. (n.d.). In Wikipedia. https://en.wikipedia.org/wiki/Exponential_backoff
