package core

import "regexp"

type info struct {
	resolutionRe *regexp.Regexp
	sampleRateRe *regexp.Regexp
	bitrateRe    *regexp.Regexp
}

func newInfo() *info {
	return &info{
		resolutionRe: regexp.MustCompile(`^\d+x\d+$`),
		sampleRateRe: regexp.MustCompile(`^\d+\s*Hz`),
		bitrateRe:    regexp.MustCompile(`^\d+\s*kb/s`),
	}
}
