package inltest

import (
	"testing"
)

func TestCheckInlineable(t *testing.T) {
	// Few symbols from the stdlib that ought to be
	// inlineable (they have inlining tests inside the compiler).

	issues, err := CheckInlineable(map[string][]string{
		"math/big": {
			"bigEndianWord",
		},
		"bytes": {
			"(*Buffer).Cap",
			"(*Buffer).Len",
		},
	})

	if err != nil {
		t.Fatalf("inltest failed: %v", err)
	}
	for _, issue := range issues {
		t.Error(issue)
	}
}
