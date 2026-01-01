package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func AuthMiddleware(cfg *viper.Viper) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil Header Authorization
		autHeader := c.Get("Authorization")
		if autHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		// 2. Format Harus "Bearer <Token>"
		parts := strings.Split(autHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization Format",
			})
		}
		tokenString := parts[1]

		// 3. Parse & Validasi Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			secret := cfg.GetString("jwt.secret")
			return []byte(secret), nil
		})

		// 4. Cek Error Parse
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or Expired token",
			})
		}

		// 5. Ekstrak Claims (Data User)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Token Claims",
			})
		}

		// 6. Simpan USer ID ke Contex (Agar bisa di pakai di controller)
		c.Locals("user_id", claims["sub"])
		c.Locals("email", claims["email"])

		return c.Next()
	}
}
