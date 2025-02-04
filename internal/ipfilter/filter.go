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
		ip := strings.TrimSpace(scanner.Text())
		if ip != "" {
			newIpMap[ip] = true
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
