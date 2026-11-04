package executor

import (
	idto "github.com/neosy/elengrab/internal/app/services/ytdlp/internal/downloader/dto"
)

func addCookiesToArgs(args []string, opts ...idto.ExecutorOption) []string {
	options := idto.NewExecutorOptions(opts...)

	// Append the cookies file path to the arguments if it exists
	if options.CookieFilePath != "" {
		args = append(args, "--cookies", options.CookieFilePath)
	}

	return args
}
