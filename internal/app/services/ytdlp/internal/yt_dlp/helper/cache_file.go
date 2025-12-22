package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func FormatCacheFileName(url string) string {
	hash := sha256.Sum256([]byte(url))
	return fmt.Sprintf("%s.json", hex.EncodeToString(hash[:16]))
}
