package handler

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strings"

	"simpleConfigServer/internal/auth"
	"simpleConfigServer/internal/config"
	"simpleConfigServer/internal/ipfilter"
	"simpleConfigServer/internal/logger"
	"simpleConfigServer/internal/rate_limiter"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	logger.Log.Printf("Endpoint hit: %s", r.URL.Path)

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	logger.Log.Printf("IP address: %s", ip)
	if err != nil {
		logger.Log.Printf("Invalid IP address: %s", r.RemoteAddr)
		http.Error(w, "Invalid IP address", http.StatusForbidden)
		return
	}

	if !ipfilter.IsIPAllowed(ip) {
		logger.Log.Printf("IP not allowed: %s", ip)
		http.Error(w, "IP not allowed", http.StatusForbidden)
		return
	}

	limiter := rate_limiter.GetRateLimiter(ip)
	if !limiter.Allow() {
		logger.Log.Printf("Rate limit exceeded for %s", ip)
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !auth.ValidateJWT(tokenString) {
		logger.Log.Printf("Unauthorized access attempt to %s", r.URL.Path)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := strings.Split(r.URL.Path, "/")
	if len(vars) < 4 {
		logger.Log.Printf("Invalid request path: %s", r.URL.Path)
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}

	product, env, configKey := vars[1], vars[2], vars[3]
	logger.Log.Printf("Product: %s, Environment: %s, Config Key: %s", product, env, configKey)

	if env != "staging" && env != "production" && env != "development" {
		logger.Log.Printf("Unsupported environment: %s", env)
		http.Error(w, "Environment not supported", http.StatusNotFound)
		return
	}

	configs := config.GetConfigs()
	productConfigs, exists := configs[product]
	if !exists {
		logger.Log.Printf("Product not found: %s", product)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	configValue, found := productConfigs[configKey]
	if !found {
		logger.Log.Printf("Configs not found: %s for product %s", configKey, product)
		http.Error(w, "Configs not found", http.StatusNotFound)
		return
	}

	logger.Log.Printf("Successfully retrieved configs: %s for product %s", configKey, product)
	response := map[string]string{configKey: configValue}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	json.NewEncoder(w).Encode(response)
}
