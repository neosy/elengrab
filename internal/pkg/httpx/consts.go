package httpx

const (
	LinuxUserAgent = "Mozilla/5.0 (X11; Linux x86_64)" +
		" AppleWebKit/537.36 (KHTML, like Gecko)" +
		" Chrome/146.0.0.0 Safari/537.36"

	acceptLanguageDefault = "en-US,en;q=0.9"

	// Limit for the size of the image to avoid excessive memory usage.
	LimitImageDefault int64 = 5 << 20 // 5 MB

	// Limit for the size of the HTML to avoid excessive memory usage.
	LimitHTMLDefault int64 = 512 << 10 // 512 KB

	// Limit for the size of the HTML for reading the TITLE to avoid excessive memory usage.
	LimitHTMLTitleDefault int64 = 1024 << 10 // 1024 KB
)
