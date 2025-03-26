package ipfilter

import (
	"bufio"
	"os"
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
	}
	defer file.Close()

	newIpMap := make(map[string]bool)

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
	}

	mu.Lock()
	allowedIPs = newIpMap
	mu.Unlock()
}

func IsIPAllowed(ip string) bool {
	mu.RLock()
	defer mu.RUnlock()
	logger.Log.Printf("Checking if IP %s is allowed", ip)
	if len(allowedIPs) == 0 {
		logger.Log.Println("Allowed IPs list is empty, allowing all IPs")
		return true
	}
	isAllowed, exists := allowedIPs[ip]
	if !exists {
		logger.Log.Printf("IP %s is not allowed", ip)
	}
	return isAllowed
}
