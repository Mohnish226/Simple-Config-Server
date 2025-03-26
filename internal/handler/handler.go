package handler

import (
	"os"
	"strings"

	"simpleConfigServer/internal/auth"
	"simpleConfigServer/internal/config"
	"simpleConfigServer/internal/ipfilter"
	"simpleConfigServer/internal/logger"
	"simpleConfigServer/internal/rate_limiter"

	"github.com/gofiber/fiber/v2"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func ConfigHandler(c *fiber.Ctx) error {
	logger.Log.Printf("Endpoint hit: %s", c.Path())

	// Get IP address from Fiber context
	ip := c.IP()
	logger.Log.Printf("IP address: %s", ip)

	// Check if IP is allowed
	if !ipfilter.IsIPAllowed(ip) {
		logger.Log.Printf("IP not allowed: %s", ip)
		return c.Status(fiber.StatusForbidden).SendString("IP not allowed")
	}

	limiter := rate_limiter.GetRateLimiter(ip)
	if !limiter.Allow() {
		logger.Log.Printf("Rate limit exceeded for %s", ip)
		return c.Status(fiber.StatusTooManyRequests).SendString("Too Many Requests")
	}

	tokenString := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if !auth.ValidateJWT(tokenString) {
		logger.Log.Printf("Unauthorized access attempt to %s", c.Path())
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	vars := strings.Split(c.Path(), "/")
	if len(vars) < 4 {
		logger.Log.Printf("Invalid request path: %s", c.Path())
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request path")
	}

	product, env, configKey := vars[1], vars[2], vars[3]
	logger.Log.Printf("Product: %s, Environment: %s, Config Key: %s", product, env, configKey)

	if env != "staging" && env != "production" && env != "development" {
		logger.Log.Printf("Unsupported environment: %s", env)
		return c.Status(fiber.StatusNotFound).SendString("Environment not supported")
	}

	configs := config.GetConfigs()
	productConfigs, exists := configs[product]
	if !exists {
		logger.Log.Printf("Product not found: %s", product)
		return c.Status(fiber.StatusNotFound).SendString("Product not found")
	}

	configValue, found := productConfigs[configKey]
	if !found {
		logger.Log.Printf("Configs not found: %s for product %s", configKey, product)
		return c.Status(fiber.StatusNotFound).SendString("Configs not found")
	}

	logger.Log.Printf("Successfully retrieved configs: %s for product %s", configKey, product)
	response := map[string]string{configKey: configValue}

	// Set security headers
	c.Set("Content-Type", "application/json")
	c.Set("Content-Security-Policy", "default-src 'self'")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-XSS-Protection", "1; mode=block")

	return c.JSON(response)
}
