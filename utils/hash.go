package utils

import (
	"crypto/sha256"
	"fmt"
)

func Sha256String(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
