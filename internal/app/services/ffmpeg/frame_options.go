package ffmpegsrv

type FrameOptions struct {
	Strategy FrameStrategy
	Format   FrameFormat
}

type FrameOption func(opts *FrameOptions)

// WithFrameStrategy sets the frame extraction strategy for GetBestFrame.
func WithFrameStrategy(strategy FrameStrategy) FrameOption {
	return func(opts *FrameOptions) {
		opts.Strategy = strategy
	}
}

// WithFrameFormat sets the frame output format for GetBestFrame.
func WithFrameFormat(format FrameFormat) FrameOption {
	return func(opts *FrameOptions) {
		opts.Format = format
	}
}
