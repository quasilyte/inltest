[![Go Report Card](https://goreportcard.com/badge/github.com/Quasilyte/inltest)](https://goreportcard.com/report/github.com/Quasilyte/inltest)
[![GoDoc](https://godoc.org/github.com/Quasilyte/inltest?status.svg)](https://godoc.org/github.com/Quasilyte/inltest)
[![Build Status](https://travis-ci.org/Quasilyte/inltest.svg?branch=master)](https://travis-ci.org/Quasilyte/inltest)

# inltest

Package inltest helps you to test that performance-sensitive funcs are inlineable.

Usually should be used inside your tests, so you can see that some functions are
not inlineable anymore due to, for example, cost increase during the last refactoring.

> Note: please don't try to interpret returned "not inlined resons" slice.
> Its contents may change from one Go version to another.
> The only information you can safely rely on is whether function
> is inlineable or not. And usually you want all functions from the
> input map to be inlineable (otherwise why would you include them)?

## Installation

```bash
go get -v github.com/Quasilyte/inltest
```

## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/Quasilyte/inltest"
)

func main() {
	issues, err := inltest.CheckInlineable(map[string][]string{
		"github.com/Quasilyte/inltest": {
			"CheckInlineable",
			"nonexisting",
		},

		// errors.New is inlineable => gives no issue.
		"errors": {
			"New",
		},

		"strings": {
			"(*Builder).WriteRune",
		},
	})
	if err != nil {
		log.Fatalf("inltest failed: %v", err)
	}
	for _, issue := range issues {
		fmt.Println(issue)
	}
}
```

For tests, you can do something like:

```go
func TestInlining(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	issues, err := CheckInlineable(map[string][]string{
		"my/important/pkg": {
			"func1",
			"func2",
			"(*Value).Set",
		},
	})
	if err != nil {
		t.Fatalf("inltest failed: %v", err)
	}
	for _, issue := range issues {
		t.Error(issue)
	}
}
```
