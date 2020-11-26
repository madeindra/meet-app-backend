package helpers

import (
	"math/rand"
	"time"
)

type RandomInterface interface {
	RandomString(length int) string
}

type RandomImplementation struct {
	random RandomInterface
}

func NewRandomHelper() *RandomImplementation {
	return &RandomImplementation{}
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func randomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (random *RandomImplementation) RandomString(length int) string {
	return randomStringWithCharset(length, charset)
}
