package utils

import "github.com/gofiber/fiber/v2"

func Cookie(key string, value string) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = key
	cookie.HTTPOnly = true
	cookie.Value = value

	if Get("ENV") == "development" {
		cookie.Secure = false
	} else {
		cookie.Domain = ".gists.app" // hardcoded
		cookie.Secure = true
	}
	return cookie
}
