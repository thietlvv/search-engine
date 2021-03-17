package hash

import (
	"crypto/sha512"
	"encoding/hex"
)

func SHA512(content string) string {
	h := sha512.New()
	h.Write([]byte(content))
	sha := h.Sum(nil)
	return hex.EncodeToString(sha)
}
