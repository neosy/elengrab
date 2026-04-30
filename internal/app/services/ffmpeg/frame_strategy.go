package ffmpegsrv

import "fmt"

// FrameStrategy defines the interface for different frame extraction strategies.
type FrameStrategy interface {
	// Args returns the ffmpeg arguments for the frame extraction strategy.
	Args() []string
}

// ThumbnailStrategy defines the strategy for extracting a thumbnail frame.
type FrameStrategyThumbnail struct {
	// BatchSize is the number of frames to analyze for thumbnail extraction. Default is 300.
	BatchSize int
}

func (s FrameStrategyThumbnail) Args() []string {
	batch := s.BatchSize
	if batch <= 0 {
		batch = 300
	}

	return []string{
		"-vf", "thumbnail=300",
		"-frames:v", "1",
	}
}

// FrameStrategySceneChange defines the strategy for extracting a frame based on scene changes.
type FrameStrategySceneChange struct {
	// Threshold is the scene change detection threshold. Default is 0.4.
	Threshold float64
}

func (s FrameStrategySceneChange) Args() []string {
	threshold := s.Threshold
	if threshold <= 0 {
		threshold = 0.4
	}

	return []string{
		"-vf", fmt.Sprintf("select=gt(scene\\,%f)", threshold),
		"-frames:v", "1",
	}
}

// FrameStrategyBalanced defines a strategy that combines scene detection
// and thumbnail selection to extract the most representative frame.
type FrameStrategyBalanced struct {
	// Threshold is the scene change detection threshold. Default is 0.3.
	Threshold float64

	// SampleSize is the number of frames passed to thumbnail filter. Default is 5.
	SampleSize int
}

func (s FrameStrategyBalanced) Args() []string {
	threshold := s.Threshold
	if threshold <= 0 {
		threshold = 0.3
	}

	sampleSize := s.SampleSize
	if sampleSize <= 0 {
		sampleSize = 5
	}

	return []string{
		"-vf", fmt.Sprintf("select=gt(scene\\,%f),thumbnail=%d", threshold, sampleSize),
		"-frames:v", "1",
	}
}
