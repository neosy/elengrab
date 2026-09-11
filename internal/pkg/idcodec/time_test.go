package idcodec

import (
	"testing"
	"time"
)

func TestEncodeTimeBase64URL(t *testing.T) {
	location := time.FixedZone("UTC+3", 3*60*60)

	value := time.Date(
		2026,
		9,
		6,
		20,
		2,
		40,
		123456789,
		location,
	)

	got := EncodeTimeBase64URL(value)
	want := EncodeInt64Base64URL(value.Unix())

	if got != want {
		t.Errorf("EncodeTimeBase64URL() = %q, want %q", got, want)
	}
}

func TestDecodeTimeBase64URL(t *testing.T) {
	value := time.Date(
		2026,
		9,
		6,
		17,
		2,
		40,
		0,
		time.UTC,
	)

	encoded := EncodeTimeBase64URL(value)

	got, err := DecodeTimeBase64URL(encoded)
	if err != nil {
		t.Fatalf("DecodeTimeBase64URL() error = %v", err)
	}

	if !got.Equal(value) {
		t.Errorf("DecodeTimeBase64URL() = %v, want %v", got, value)
	}

	if got.Location() != time.UTC {
		t.Errorf("DecodeTimeBase64URL() location = %v, want UTC", got.Location())
	}
}

func TestDecodeTimeBase64URLInvalid(t *testing.T) {
	_, err := DecodeTimeBase64URL("invalid!")

	if err == nil {
		t.Fatal("DecodeTimeBase64URL() error = nil, want error")
	}
}

func TestEncodeDecodeTimeBase64URL(t *testing.T) {
	values := []time.Time{
		time.Unix(0, 0).UTC(),
		time.Unix(1788706800, 0).UTC(),
		time.Date(
			2026,
			9,
			6,
			20,
			2,
			40,
			0,
			time.FixedZone("UTC+3", 3*60*60),
		),
	}

	for _, value := range values {
		t.Run("round trip", func(t *testing.T) {
			encoded := EncodeTimeBase64URL(value)

			got, err := DecodeTimeBase64URL(encoded)
			if err != nil {
				t.Fatalf("DecodeTimeBase64URL() error = %v", err)
			}

			if !got.Equal(value) {
				t.Errorf("round trip = %v, want %v", got, value)
			}

			if got.Location() != time.UTC {
				t.Errorf("round trip location = %v, want UTC", got.Location())
			}
		})
	}
}

func TestEncodeTimeBase64URLTruncatesNanoseconds(t *testing.T) {
	value := time.Date(
		2026,
		9,
		6,
		17,
		2,
		40,
		123456789,
		time.UTC,
	)

	encoded := EncodeTimeBase64URL(value)

	got, err := DecodeTimeBase64URL(encoded)
	if err != nil {
		t.Fatalf("DecodeTimeBase64URL() error = %v", err)
	}

	want := value.Truncate(time.Second)

	if !got.Equal(want) {
		t.Errorf("DecodeTimeBase64URL() = %v, want %v", got, want)
	}
}
