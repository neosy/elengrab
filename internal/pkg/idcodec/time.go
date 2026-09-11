package idcodec

import "time"

// EncodeTimeBase64URL converts time.Time into a URL-safe short string representation.
func EncodeTimeBase64URL(value time.Time) string {
	return EncodeInt64Base64URL(value.Unix())
}

// DecodeTimeBase64URL converts a URL-safe short string representation into time.Time.
func DecodeTimeBase64URL(value string) (time.Time, error) {
	timestamp, err := DecodeInt64Base64URL(value)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(timestamp, 0).UTC(), nil
}
