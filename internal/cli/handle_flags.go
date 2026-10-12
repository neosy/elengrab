package cli

import (
	"fmt"

	iconfig "github.com/neosy/elengrab/internal/config"
	"github.com/spf13/pflag"
)

// Handle handles command-line flags that don't require application startup.
func (flags Flags) Handle() error {
	if flags.Help {
		pflag.Usage()
	}

	if flags.Version {
		fmt.Printf("%s %s\n", iconfig.AppName, iconfig.AppVersion)
	}

	return nil
}
