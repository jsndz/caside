package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)
var secretKey = []byte("jwtSecret")

func Authenticate(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")
	token ,err := jwt.Parse(tokenString,func (t *jwt.Token)(interface{},error){
		return secretKey,nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Locals("id", claims["sub"])
	}
	return c.Next()
}