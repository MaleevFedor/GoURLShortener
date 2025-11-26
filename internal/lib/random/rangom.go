package random

import (
	"math/rand"
	"time"
)

var rnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rnd.Intn(len(letters))]
	}
	return string(b)
}
