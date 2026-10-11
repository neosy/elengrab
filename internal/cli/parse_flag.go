package cli

import (
	"fmt"
	"os"

	iconfig "github.com/neosy/elengrab/internal/config"
	"github.com/spf13/pflag"
)

// ParseFlag parses command-line flags and handles the version flag.
func ParseFlags() {
	versionFlag := pflag.BoolP("version", "v", false, "Show version")

	pflag.Parse()

	if versionFlag != nil && *versionFlag {
		fmt.Printf("%s - video downloader with format and quality options\n", iconfig.AppName)
		fmt.Printf("Version %s\n", iconfig.AppVersion)
		os.Exit(0)
	}
}
