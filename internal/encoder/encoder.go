package encoder

import (
	"strings"
)

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func EncodeBase62(id uint64) string {
	if id == 0 {
		return "0"
	}

	var encoded strings.Builder
	for id > 0 {
		remainder := id % 62
		encoded.WriteByte(base62Chars[remainder])
		id /= 62
	}

	// Reverse the encoded string
	runes := []rune(encoded.String())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func DecodeBase62(encoded string) uint64 {
	var decoded uint64
	for _, c := range encoded {
		decoded *= 62
		switch {
		case c >= '0' && c <= '9':
			decoded += uint64(c - '0')
		case c >= 'a' && c <= 'z':
			decoded += uint64(c - 'a' + 10)
		case c >= 'A' && c <= 'Z':
			decoded += uint64(c - 'A' + 36)
		}
	}
	return decoded
}
