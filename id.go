package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

var b string

func sessID(a byte) string {
	b := make([]byte, a)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)

}
