package iconfig

import (
	"fmt"
	"os"

	iconstants "github.com/neosy/elengrab/infrastructure/constants"
	"github.com/spf13/pflag"
)

func parseFlag() {
	versionFlag := pflag.BoolP("version", "v", false, "Show version")

	pflag.Parse()

	if versionFlag != nil && *versionFlag {
		fmt.Printf("%s - fast YouTube video/audio grabber with format and quality options\n", iconstants.AppName)
		fmt.Printf("Version %s\n", iconstants.AppVersion)
		os.Exit(0)
	}
}
