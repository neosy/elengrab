package clientcap

func isLegacyIOS(major int) bool {
	// iOS 12 and below is a problem area for modern JS/HTMX 2.x
	return major <= 12
}
