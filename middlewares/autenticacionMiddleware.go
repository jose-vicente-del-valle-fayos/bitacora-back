package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"nd-back/utilidades"
)

// Autenticado tests if some user is authenticated by parsing a JWT
func Autenticado(c *fiber.Ctx) error {
	cookie := c.Cookies("btcr-jwt")
	if _, err := utilidades.ParsearJWT(cookie); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"mensaje": "sin autenticar",
		})
	}
	return c.Next()
}
