package clientcap

// Capabilities describes the detected features and limitations of the client
// based on the User-Agent string. It is used to enable adaptive delivery of
// frontend assets and behavior depending on browser and OS capabilities.
type Capabilities struct {
	IsIOS    bool
	IOSMajor int
	IOSMinor int
	IOSPatch int

	IsLegacyWebKit bool
}
