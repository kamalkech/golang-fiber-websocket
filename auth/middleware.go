package auth

import (
	jwtware "github.com/gofiber/contrib/jwt"
)

var JwtMiddleware = jwtware.New(jwtware.Config{
	SigningKey: jwtware.SigningKey{
		Key: []byte("secret"),
	},
})
