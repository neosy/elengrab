package menu

import dlink "github.com/neosy/elengrab/internal/domain/link"

type MenuActionOptions struct {
	ErrorText    string
	HasMetadata  bool
	HasShareLink bool
}

type MenuActionOption func(*MenuActionOptions)

func applyMenuActionOptions(options *MenuActionOptions, opts ...MenuActionOption) {
	for _, opt := range opts {
		opt(options)
	}
}

func NewMenuActionOptions(opts ...MenuActionOption) MenuActionOptions {
	options := MenuActionOptions{}

	applyMenuActionOptions(&options, opts...)

	return options
}

func WithErrorText(text string) MenuActionOption {
	return func(m *MenuActionOptions) {
		m.ErrorText = text
	}
}

func WithShareLink(link *dlink.Link) MenuActionOption {
	return func(m *MenuActionOptions) {
		m.HasShareLink = link != nil
	}
}

func WithMetadata(metadataAvailable bool) MenuActionOption {
	return func(m *MenuActionOptions) {
		m.HasMetadata = metadataAvailable
	}
}
