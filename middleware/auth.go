package middleware

import (
	"fibox/jwt"
	"fibox/response"
	"strings"

	"github.com/gofiber/fiber/v3"
)

type AuthInfo struct {
	UserID string `json:"userID"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

func AuthMiddleware(jwtSvc *jwt.JWTService) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Authorization header is required")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Unauthorized(c, "Invalid authorization header format")
		}

		tokenString := parts[1]
		claims, err := jwtSvc.ValidateToken(tokenString)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				return response.Unauthorized(c, "Token has expired")
			}
			return response.Unauthorized(c, "Invalid token")
		}

		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func GetAuthInfo(c fiber.Ctx) AuthInfo {
	return AuthInfo{
		UserID: c.Locals("userID").(string),
		Email:  c.Locals("email").(string),
		Role:   c.Locals("role").(string),
	}
}
