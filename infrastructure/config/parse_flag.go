package iconfig

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

func parseFlag() {
	versionFlag := pflag.BoolP("version", "v", false, "Show version")

	pflag.Parse()

	if versionFlag != nil && *versionFlag {
		fmt.Printf("%s - fast YouTube video/audio grabber with format and quality options\n", AppName)
		fmt.Printf("Version %s\n", AppVersion)
		os.Exit(0)
	}
}
