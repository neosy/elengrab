package iconfetcher

type iconCandidate struct {
	url     string
	size    int // max(width, height), 0 if unknown
	rel     string
	imgType string
}
