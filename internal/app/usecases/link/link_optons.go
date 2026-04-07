package link

const defaultShortCodeLength = 8

type LinkOption func(*LinkOptions)

type LinkOptions struct {
	// Base URL for short links
	BaseURL string
	// Length of the generated short code
	ShortCodeLength uint8
	//
	Deterministic bool
}

func DefaultLinkOptions() LinkOptions {
	return LinkOptions{
		ShortCodeLength: defaultShortCodeLength,
		Deterministic:   false,
	}
}

func LinkOptionBaseURL(url string) LinkOption {
	return func(o *LinkOptions) {
		o.BaseURL = url
	}
}

func LinkOptionShortCodeLength(length uint8) LinkOption {
	return func(o *LinkOptions) {
		o.ShortCodeLength = length
	}
}

func LinkOptionDeterministic(deterministic bool) LinkOption {
	return func(o *LinkOptions) {
		o.Deterministic = deterministic
	}
}
