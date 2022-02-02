# retry

[![.github/workflows/test.yaml](https://github.com/kei6u/retry/actions/workflows/test.yaml/badge.svg)](https://github.com/kei6u/retry/actions/workflows/test.yaml)
[![GoDoc](https://godoc.org/github.com/kei6u/retry?status.svg&style=flat-square)](http://godoc.org/github.com/kei6u/retry)

This is a Go library provides retry functionality for general operations such as constant interval retry and exponential backoff algorithms.

## Motivation

There are popular awesome similar packages; [avast/retry-go](https://github.com/avast/retry-go), [lestrrat-go/backoff](https://github.com/lestrrat-go/backoff).
However, I want a new package provides more simple interface and implementation. This is a biggest and only motivation to create this package.

## Usage

### Import

```bash
go get github.com/kei6u/retry
```

```go
import "github.com/kei6u/retry"
```

### Constant

The constant retry retries at a constant intervals.

### Exponential backoff

"Exponential backoff is an algorithm that uses feedback to multiplicatively decrease the rate of some process, in order to gradually find an acceptable rate. These algorithms find usage in a wide range of systems and processes, with radio networks and computer networks being particularly notable." ("Exponential backoff," n.d.)

Reference

Exponential backoff. (n.d.). In Wikipedia. https://en.wikipedia.org/wiki/Exponential_backoff
