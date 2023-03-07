package main

import "math/rand"

// RandomChar generate  pseudo random chars
func RandomChar(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	str := make([]rune, n)
	for i := range str {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}
