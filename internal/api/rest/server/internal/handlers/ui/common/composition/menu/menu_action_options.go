package menu

import dlink "github.com/neosy/elengrab/internal/domain/link"

type MenuActionOptions struct {
	ErrorText    string
	HasShareLink bool
}

type MenuActionOption func(*MenuActionOptions)

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
