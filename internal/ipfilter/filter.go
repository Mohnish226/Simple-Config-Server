package ipfilter

import (
	"bufio"
	"os"
	"simpleConfigServer/internal/audit"
	"simpleConfigServer/internal/logger"
	"strings"
	"sync"
)

var (
	allowedIPs = make(map[string]bool)
	mu         sync.RWMutex
)

func LoadAllowedIPs(AllowedIPsFile string) {
	file, err := os.Open(AllowedIPsFile)
	if err != nil {
		logger.Log.Fatalf("Failed to open allowed IPs file: %v", err)
		audit.LogSystem("IP_FILTER_LOAD", "FAILED", map[string]interface{}{
			"file":  AllowedIPsFile,
			"error": err.Error(),
		})
	}
	defer file.Close()

	newIpMap := make(map[string]bool)
	oldIpMap := make(map[string]bool)
	mu.RLock()
	for ip := range allowedIPs {
		oldIpMap[ip] = true
	}
	mu.RUnlock()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Remove any inline comments
		if idx := strings.Index(line, "#"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		if line != "" {
			newIpMap[line] = true
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Log.Fatalf("Error reading allowed IPs file: %v", err)
		audit.LogSystem("IP_FILTER_LOAD", "FAILED", map[string]interface{}{
			"file":  AllowedIPsFile,
			"error": err.Error(),
		})
	}

	mu.Lock()
	allowedIPs = newIpMap
	mu.Unlock()

	// Log IP changes
	for ip := range newIpMap {
		if _, exists := oldIpMap[ip]; !exists {
			audit.LogSystem("IP_FILTER_CHANGE", "ADDED", map[string]interface{}{
				"ip": ip,
			})
		}
	}

	for ip := range oldIpMap {
		if _, exists := newIpMap[ip]; !exists {
			audit.LogSystem("IP_FILTER_CHANGE", "REMOVED", map[string]interface{}{
				"ip": ip,
			})
		}
	}

	audit.LogSystem("IP_FILTER_LOAD", "SUCCESS", map[string]interface{}{
		"file": AllowedIPsFile,
	})
}

func IsIPAllowed(ip string) bool {
	mu.RLock()
	defer mu.RUnlock()

	if len(allowedIPs) == 0 {
		logger.Log.Println("Allowed IPs list is empty, allowing all IPs")
		audit.LogSecurity(ip, "ALLOWED", "IP_FILTER", map[string]interface{}{
			"reason": "IP list empty",
		})
		return true
	}

	isAllowed, exists := allowedIPs[ip]
	if !exists {
		logger.Log.Printf("IP %s is not allowed", ip)
		audit.LogSecurity(ip, "DENIED", "IP_FILTER", map[string]interface{}{
			"reason": "IP not in allowed list",
		})
	} else {
		audit.LogSecurity(ip, "ALLOWED", "IP_FILTER", map[string]interface{}{
			"reason": "IP in allowed list",
		})
	}
	return isAllowed
}
