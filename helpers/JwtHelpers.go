package helpers

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func VerifyToken(c *fiber.Ctx) error {
	splitToken := strings.Split(c.Get("Authorization"), "Bearer ")
	tokenString := splitToken[1]
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		res := GetResponse(fiber.StatusUnauthorized, nil, err)
		return c.Status(res.Status).JSON(res)
	}
	c.Locals("Username", claims["username"].(string))
	c.Locals("UserId", claims["userid"])
	c.Locals("Email", claims["email"].(string))
	//c.Locals("RoleId", claims["roleid"].(string))
	return c.Next()
}
