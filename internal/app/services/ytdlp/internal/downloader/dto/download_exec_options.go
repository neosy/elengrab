package idto

import uptr "github.com/neosy/elengrab/internal/pkg/utils/pointer"

type DownloadExecOptions struct {
	ConcurrentFragments uint8
	CookieFilePath      string
	Extractor           string
	ExtractorArgs       *string
	Args                []string
}

func (o *DownloadExecOptions) Copy() *DownloadExecOptions {
	if o == nil {
		return nil
	}

	options := *o
	options.ExtractorArgs = uptr.Copy(options.ExtractorArgs)
	options.Args = o.CopyArgs()

	return &options
}

func (o *DownloadExecOptions) CopyArgs() []string {
	if o == nil {
		return nil
	}

	if o.Args == nil {
		return nil
	}

	args := make([]string, len(o.Args))
	copy(args, o.Args)

	return args
}
