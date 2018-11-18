// Package inltest helps you to test that performance-sensitive funcs are inlineable.
package inltest

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"sort"
	"strings"
)

// CheckInlineable tries to find whether provided candidates are inlineable.
//
// Returns a list of reasons for symbols that were not proven to be
// inlineable. If it's nil, then all symbols in candidates map are inlineable.
// The returned string slice is sorted.
//
// The candidates maps import path to the symbols that should be checked.
//
// Here are some examples:
//	"io/ioutil": {"ReadAll"}       // Check bytes.ReadAll function
//	"bytes":     {"(*Buffer).Len"} // Check bytes.Buffer Len method
//
// Note that you can check several symbols from the same package.
func CheckInlineable(candidates map[string][]string) ([]string, error) {
	if _, ok := candidates[""]; ok {
		return nil, fmt.Errorf("empty import path is not allowed")
	}

	// The implementation is borrowed from the cmd/compile/internal/gc/inl_test.go.

	notInlinedReason := make(map[string]string)
	pkgs := make([]string, 0, len(candidates))
	for importPath, syms := range candidates {
		pkgs = append(pkgs, importPath)
		for _, sym := range syms {
			fullName := importPath + "." + sym
			if _, ok := notInlinedReason[fullName]; ok {
				return nil, fmt.Errorf("duplicate func: %s", fullName)
			}
			notInlinedReason[fullName] = "unknown reason"
		}
	}

	args := append([]string{"build", "-a", "-gcflags=all=-m -m"}, pkgs...)
	cmd := exec.Command("go", args...)
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw
	cmdErr := make(chan error, 1)
	go func() {
		cmdErr <- cmd.Run()
		pw.Close()
	}()
	scanner := bufio.NewScanner(pr)
	curPkg := ""
	canInline := regexp.MustCompile(`: can inline ([^ ]*)`)
	haveInlined := regexp.MustCompile(`: inlining call to ([^ ]*)`)
	cannotInline := regexp.MustCompile(`: cannot inline ([^ ]*): (.*)`)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "# ") {
			curPkg = line[2:]
			continue
		}
		if m := haveInlined.FindStringSubmatch(line); m != nil {
			sym := m[1]
			delete(notInlinedReason, curPkg+"."+sym)
			continue
		}
		if m := canInline.FindStringSubmatch(line); m != nil {
			sym := m[1]
			fullName := curPkg + "." + sym
			delete(notInlinedReason, fullName)
			continue
		}
		if m := cannotInline.FindStringSubmatch(line); m != nil {
			sym, reason := m[1], m[2]
			fullName := curPkg + "." + sym
			if _, ok := notInlinedReason[fullName]; ok {
				notInlinedReason[fullName] = reason
			}
			continue
		}
	}

	issues := make([]string, 0, len(notInlinedReason))
	for fullName, reason := range notInlinedReason {
		issues = append(issues, fullName+": "+reason)
	}
	sort.Strings(issues)
	return issues, nil
}
