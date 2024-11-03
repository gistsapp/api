package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func Cookie(key string, value string) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = key
	cookie.HTTPOnly = true
	cookie.Value = value
	cookie.Expires = time.Now().Add(time.Hour * 24 * 30 * 12) // 1 year

	if Get("ENV") == "development" {
		cookie.Secure = false
	} else {
		cookie.Domain = ".gists.app" // hardcoded
		cookie.Secure = true
	}
	return cookie
}

func ClearCookie(key string) *fiber.Cookie {
	cookie := new(fiber.Cookie)
	cookie.Name = key
	cookie.HTTPOnly = true
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Hour)
	cookie.Secure = true
	if Get("ENV") == "development" {
		cookie.Secure = false
	} else {
		cookie.Domain = ".gists.app" // hardcoded
	}

	return cookie
}
