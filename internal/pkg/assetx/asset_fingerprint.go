package assetx

import (
	"github.com/cespare/xxhash/v2"
)

const (
	// 4MB: maximum amount of data used for hashing to avoid excessive CPU/memory usage
	maxAssetHashSizeDefault = 4 << 20

	hexDigits = "0123456789abcdef"
)

// AssetFingerprint64 computes a 64-bit xxhash fingerprint of the asset.
// If the input is larger than maxHashSize, it is truncated to limit processing cost.
//
// Note: This is NOT a cryptographic hash. It is intended for fast cache-busting
// and asset versioning (e.g. URLs, ETag-like identifiers).
func AssetFingerprint64(data []byte) uint64 {
	if len(data) > maxAssetHashSizeDefault {
		data = data[:maxAssetHashSizeDefault]
	}

	return xxhash.Sum64(data)
}

// AssetFingerprint32 computes a 32-bit fingerprint derived from the 64-bit hash.
//
// It is a truncated version of the 64-bit fingerprint and has higher collision risk.
// Suitable only for cases where compactness is preferred over uniqueness (e.g. short URLs).
func AssetFingerprint32(data []byte) uint32 {
	h := AssetFingerprint64(data)
	return uint32(h & 0xffffffff)
}

// AssetFingerprintHex returns a 64-bit fingerprint encoded as a hex string.
//
// This is the standard representation for use in URLs, cache keys, and ETag-like headers.
func AssetFingerprintHex(data []byte) string {
	h := AssetFingerprint64(data)

	var b [16]byte
	for i := 15; i >= 0; i-- {
		b[i] = hexDigits[h&0xf]
		h >>= 4
	}

	return string(b[:])
}

// AssetFingerprintHex32 returns a 32-bit fingerprint (derived from 64-bit hash)
// encoded as a hex string.
//
// This produces a shorter identifier but with significantly higher collision probability.
// Intended only for non-critical caching or compact URL representations.
func AssetFingerprintHex32(data []byte) string {
	h := AssetFingerprint32(data)

	var b [8]byte
	for i := 7; i >= 0; i-- {
		b[i] = hexDigits[h&0xf]
		h >>= 4
	}

	return string(b[:])
}
