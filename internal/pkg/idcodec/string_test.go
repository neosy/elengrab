package idcodec

import "testing"

func TestBase64URLRoundTrip(t *testing.T) {
	tests := []string{
		"",
		"viewMode=new&lastId=5E7XaHL5Tx6EPmvOM_IABA&lastCreateAt=2026-09-06T17:02:40Z&lastViewsKey=1",
		"viewMode=popular&lastId=5E7XaHL5Tx6EPmvOM_IABA&lastCreateAt=2026-09-06T17:02:40Z&lastViewsKey=123",
		"viewMode=old&lastId=5E7XaHL5Tx6EPmvOM_IABA&lastCreateAt=2026-09-06T17:02:40Z",
	}

	for _, value := range tests {
		t.Run(value, func(t *testing.T) {
			encoded := EncodeStringBase64URL(value)

			got, err := DecodeStringBase64URL(encoded)
			if err != nil {
				t.Fatalf("DecodeStringBase64URL() error = %v", err)
			}

			if got != value {
				t.Fatalf("round trip = %q, want %q", got, value)
			}
		})
	}
}

func TestEncodeStringBase64URL(t *testing.T) {
	value := "viewMode=new&lastId=5E7XaHL5Tx6EPmvOM_IJABA&lastCreateAt=2026-09-06T17:02:40Z&lastViewsKey=1"
	want := "dmlld01vZGU9bmV3Jmxhc3RJZD01RTdYYUhMNVR4NkVQbXZPTV9JSkFCQSZsYXN0Q3JlYXRlQXQ9MjAyNi0wOS0wNlQxNzowMjo0MFombGFzdFZpZXdzS2V5PTE"

	got := EncodeStringBase64URL(value)

	if got != want {
		t.Fatalf("EncodeStringBase64URL() = %q, want %q", got, want)
	}
}
