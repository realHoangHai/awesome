package health

import (
	"context"
)

type (
	// CheckFunc is quick way to define a health checker.
	CheckFunc func(context.Context) error

	// Checker provides functionality to check health of the service.
	Checker interface {
		// CheckHealth setup health checke to the target service.
		// Return error if the target service is not available for working.
		CheckHealth(context.Context) error
	}
)

var (
	_ Checker = (CheckFunc)(nil)
)

// CheckHealth implements Checker interface.
func (cf CheckFunc) CheckHealth(ctx context.Context) error {
	return cf(ctx)
}
