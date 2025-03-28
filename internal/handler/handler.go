package handler

import (
	"os"
	"strings"

	"simpleConfigServer/internal/audit"
	"simpleConfigServer/internal/auth"
	"simpleConfigServer/internal/config"
	"simpleConfigServer/internal/ipfilter"
	"simpleConfigServer/internal/rate_limiter"

	"github.com/gofiber/fiber/v2"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func ConfigHandler(c *fiber.Ctx) error {
	ip := c.IP()
	audit.LogSystem("REQUEST", "START", map[string]interface{}{
		"path": c.Path(),
		"ip":   ip,
	})

	// Check if IP is allowed
	if !ipfilter.IsIPAllowed(ip) {
		audit.LogSecurity(ip, "DENIED", "IP_FILTER", map[string]interface{}{
			"reason": "IP not in allowed list",
		})
		return c.Status(fiber.StatusForbidden).SendString("IP not allowed")
	}

	limiter := rate_limiter.GetRateLimiter(ip)
	if !limiter.Allow() {
		audit.LogSecurity(ip, "DENIED", "RATE_LIMIT", map[string]interface{}{
			"reason": "Rate limit exceeded",
		})
		return c.Status(fiber.StatusTooManyRequests).SendString("Too Many Requests")
	}

	tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	claims, isValid := auth.ValidateJWT(tokenString)
	if !isValid {
		audit.LogAuth(ip, "FAILED", "")
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	audit.LogAuth(ip, "SUCCESS", claims.UserID)

	vars := strings.Split(c.Path(), "/")
	if len(vars) < 4 {
		audit.LogSystem("REQUEST", "INVALID", map[string]interface{}{
			"reason": "Invalid request path",
			"path":   c.Path(),
		})
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request path")
	}

	product, env, configKey := vars[1], vars[2], vars[3]

	if env != "staging" && env != "production" && env != "development" {
		audit.LogConfigAccess(ip, "DENIED", product, env, configKey, claims.UserID)
		return c.Status(fiber.StatusNotFound).SendString("Environment not supported")
	}

	configs := config.GetConfigs()
	productConfigs, exists := configs[product]
	if !exists {
		audit.LogConfigAccess(ip, "DENIED", product, env, configKey, claims.UserID)
		return c.Status(fiber.StatusNotFound).SendString("Product not found")
	}

	configValue, found := productConfigs[configKey]
	if !found {
		audit.LogConfigAccess(ip, "DENIED", product, env, configKey, claims.UserID)
		return c.Status(fiber.StatusNotFound).SendString("Configs not found")
	}

	audit.LogConfigAccess(ip, "SUCCESS", product, env, configKey, claims.UserID)
	response := map[string]string{configKey: configValue}

	// Set security headers
	c.Set("Content-Type", "application/json")
	c.Set("Content-Security-Policy", "default-src 'self'")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-XSS-Protection", "1; mode=block")

	return c.JSON(response)
}
