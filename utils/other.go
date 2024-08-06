package utils

import "math/rand"

func GenToken(length int) string {
	chars := "0123456"
	result := ""
	for i := 0; i < length; i++ {
		result += string(chars[rand.Intn(len(chars))])
	}
	return result
}
