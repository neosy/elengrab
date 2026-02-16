package ffmpeg

import "regexp"


type Info struct {
	resolutionRe *regexp.Regexp
	sampleRateRe *regexp.Regexp
	bitrateRe    *regexp.Regexp
}

func NewInfo() *Info {
	return &Info{
		resolutionRe: regexp.MustCompile(`^\d+x\d+$`),
		sampleRateRe: regexp.MustCompile(`^\d+\s*Hz`),
		bitrateRe:    regexp.MustCompile(`^\d+\s*kb/s`),
	}

}
