package random_string

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(lenght int, charset string) string {
	b := make([]byte, lenght)
	for i := range b {
		b[i] = charset[seededRand.Intn(len([]rune(charset)))]
	}
	return string(b)
}

func String(lenght int) string {
	return StringWithCharset(lenght, charset)
}
