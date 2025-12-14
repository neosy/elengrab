package iutils

import (
	"fmt"

	dtypes "github.com/neosy/elengrab/internal/domain/types"
)

// ScaleValue builds an ffmpeg scale expression based on source dimensions
// and target resolution, preserving aspect ratio.
//
// width, height — source dimensions (may be rotated)
// toResolution  — target resolution (Width/Height)
func ScaleValue(width, height uint16, toResolution dtypes.VideoResolution) string {
	// Target dimensions (e.g. 1280x720)
	targetW := toResolution.Width()
	targetH := toResolution.Height()

	// Determine source orientation:
	// true  -> landscape (width >= height)
	// false -> portrait (height > width)
	sourceIsLandscape := width >= height

	// Normalize source dimensions so srcW >= srcH for comparison
	srcW, srcH := width, height
	if !sourceIsLandscape {
		// swap for portrait to compare logically as W/H
		srcW, srcH = height, width
	}

	// If any dimension is zero, return empty (no scale)
	if targetW == 0 || targetH == 0 || srcW == 0 || srcH == 0 {
		return ""
	}

	// If source is already within target bounds, no scaling needed
	if srcW <= targetW || srcH <= targetH {
		return ""
	}

	// Build scale value for ffmpeg scale filter:
	// - landscape: limit height (targetH) -> "-1:targetH"
	// - portrait:  limit width  (targetH per your convention) -> "targetH:-1"
	// Note: we use -1 so ffmpeg recalculates the other dimension preserving aspect ratio.
	var scaleValue string
	if sourceIsLandscape {
		// landscape -> fix height
		scaleValue = fmt.Sprintf("-1:%d", targetH)
	} else {
		// portrait  -> fix width (using targetH per your convention)
		scaleValue = fmt.Sprintf("%d:-1", targetH)
	}

	return scaleValue
}
