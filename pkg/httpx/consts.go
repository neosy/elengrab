package httpx

const (
	LinuxUserAgent = "Mozilla/5.0 (X11; Linux x86_64)" +
		" AppleWebKit/537.36 (KHTML, like Gecko)" +
		" Chrome/143.0.0.0 Safari/537.36"

	acceptLanguageDefault = "en-US,en;q=0.9"

	LimitImageDefault int64 = 5 << 20   // 5 MB
	LimitHTMLDefault  int64 = 512 << 10 // 512 KB
)
