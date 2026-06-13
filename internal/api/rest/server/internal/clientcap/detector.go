package clientcap

// Detect parses the User-Agent string and returns a structured set of
// client capabilities. It is primarily used to determine whether the
// client is running a legacy iOS WebKit (e.g. iOS <= 12) in order to
// serve compatible frontend assets.
func Detect(userAgent string) Capabilities {
	major, minor, patch, ok := parseIOSVersion(userAgent)

	if !ok {
		return Capabilities{
			IsIOS: false,
		}
	}

	return Capabilities{
		IsIOS:          true,
		IOSMajor:       major,
		IOSMinor:       minor,
		IOSPatch:       patch,
		IsLegacyWebKit: isLegacyIOS(major),
	}
}
