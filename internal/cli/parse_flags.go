package cli

import (
	"github.com/spf13/pflag"
)

// ParseFlag parses command-line flags and handles the version flag.
func ParseFlags() Flags {
	var flags Flags

	pflag.CommandLine.SortFlags = false

	pflag.BoolVarP(&flags.Version, "version", "v", false, "Show version")
	pflag.BoolVarP(&flags.Help, "help", "h", false, "Help")

	pflag.Parse()

	return flags
}
