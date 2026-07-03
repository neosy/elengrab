package menu

type MenuActionOptions struct {
	ErrorText string
}

type MenuActionOption func(*MenuActionOptions)

func WithErrorText(text string) MenuActionOption {
	return func(m *MenuActionOptions) {
		m.ErrorText = text
	}
}
