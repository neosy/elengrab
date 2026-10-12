package cli

type Flags struct {
	Version bool
	Help    bool
}

// ShouldExit reports whether the application should exit after handling flags.
func (flags Flags) ShouldExit() bool {
	return flags.Help || flags.Version
}
