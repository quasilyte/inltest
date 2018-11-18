package inltest_test

import (
	"fmt"
	"log"

	"github.com/Quasilyte/inltest"
)

func Example() {
	issues, err := inltest.CheckInlineable(map[string][]string{
		"github.com/Quasilyte/inltest": {
			"CheckInlineable",
			"nonexisting",
		},

		// errors.New is inlineable => gives no issue.
		"errors": {
			"New",
		},
	})
	if err != nil {
		log.Fatalf("inltest failed: %v", err)
	}
	for _, issue := range issues {
		fmt.Println(issue)
	}

	// Output:
	// github.com/Quasilyte/inltest.CheckInlineable: unhandled op RANGE
	// github.com/Quasilyte/inltest.nonexisting: unknown reason
}
