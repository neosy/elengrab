package httpx

// MethodMethodGetOptions defines optional parameters for Get methods.
type MethodGetOptions struct {
	// Limit specifies the maximum number of bytes to read from the response.
	Limit int64
	// AcceptLanguage sets the value for the "Accept-Language" HTTP header.
	// Example: "en-US,en;q=0.9"
	AcceptLanguage string
	// IgnoreStatusCode — ignore HTTP 4xx/5xx and return the response body.
	IgnoreStatusCode bool
}
