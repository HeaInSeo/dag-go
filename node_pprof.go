//go:build pprof

package dag_go

import (
	"context"
	"runtime/pprof"
	"strconv"
)

// applyPreflightLabels attaches pprof goroutine labels to the calling goroutine.
// Called once per preFlight worker; only active when built with -tags pprof
// and DagConfig.EnablePprofLabels == true.
func applyPreflightLabels(ctx context.Context, enabled bool, nodeID string, k int) {
	if !enabled {
		return
	}
	pprof.SetGoroutineLabels(pprof.WithLabels(ctx,
		pprof.Labels("phase", "preFlight", "nodeId", nodeID, "channelIndex", strconv.Itoa(k)),
	))
}
