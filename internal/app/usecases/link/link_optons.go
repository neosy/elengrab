package link

const defaultShortCodeLength = 8

type LinkOption func(*LinkOptions)

type LinkOptions struct {
	// Base URL for short links
	BaseURL string
	// Length of the generated short code
	ShortCodeLength uint8
	//
	Deduplicate bool
}

func DefaultLinkOptions() LinkOptions {
	return LinkOptions{
		ShortCodeLength: defaultShortCodeLength,
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

func LinkOptionDeduplicate(deduplicate bool) LinkOption {
	return func(o *LinkOptions) {
		o.Deduplicate = deduplicate
	}
}
