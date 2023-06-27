package main

import (
	"crypto/sha256"
	"crypto/subtle"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

var accessKeyHash [32]byte

func keyValidator(c *fiber.Ctx, k string) (bool, error) {
	hashed := sha256.Sum256([]byte(k))

	if subtle.ConstantTimeCompare(accessKeyHash[:], hashed[:]) == 1 {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}
