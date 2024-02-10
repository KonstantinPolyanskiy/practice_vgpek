package rndutils

import (
	"math/rand"
	"strconv"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandString возвращает случайную англоязычную строку длинной n БАЙТ
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandNumberString(n int) string {
	rn := rand.Intn(n)
	return strconv.Itoa(rn)
}
