package stack_trace_error

import (
	"fmt"
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func GetStackTraceFromError(err error) []string {
	var stacks []string
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			stack := fmt.Sprintf("%+s:%d\n", f, f)
			stacks = append(stacks, stack)
		}
	}
	return stacks
}
