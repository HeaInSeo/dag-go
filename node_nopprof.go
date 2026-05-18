//go:build !pprof

package dag_go

import "context"

// applyPreflightLabels is a no-op in non-pprof builds.
// The compiler inlines and eliminates this call entirely — zero allocation, zero branch.
func applyPreflightLabels(_ context.Context, _ bool, _ string, _ int) {}
