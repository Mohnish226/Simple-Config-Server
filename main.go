package main

import (
	"flag"
	"os"
	"path/filepath"
	"simpleConfigServer/internal/config"
	"simpleConfigServer/internal/handler"
	"simpleConfigServer/internal/ipfilter"
	applogger "simpleConfigServer/internal/logger"
	"simpleConfigServer/internal/scaffolding"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Get the working directory
func getWorkingDir() string {
	execDir, err := os.Getwd()
	if err != nil {
		applogger.Log.Fatal(err)
	}
	return execDir
}

// Get configuration directory
func getConfigDir() string {
	// Check CLI flag first
	configDir := flag.String("config-dir", "", "Directory containing configuration files")
	flag.Parse()

	if *configDir != "" {
		return *configDir
	}

	// Then check environment variable
	if dir := os.Getenv("CONFIG_DIR"); dir != "" {
		return dir
	}

	// Finally, use default in current directory
	return filepath.Join(getWorkingDir(), "configurations")
}

// Get allowed IPs file path
func getAllowedIPsFile() string {
	// Check CLI flag first
	allowedIPsFile := flag.String("allowed-ips", "", "File containing allowed IP addresses")
	flag.Parse()

	if *allowedIPsFile != "" {
		return *allowedIPsFile
	}

	// Then check environment variable
	if file := os.Getenv("ALLOWED_IPS_FILE"); file != "" {
		return file
	}

	// Finally, use default in current directory
	return filepath.Join(getWorkingDir(), "allowed_ips.txt")
}

var port = func() string {
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":8080"
}()

func main() {
	// Get configuration paths
	configDir := getConfigDir()
	allowedIPsFile := getAllowedIPsFile()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Simple Config Server",
	})

	// Add middleware
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	// Setup directories and files
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		applogger.Log.Printf("Configurations directory %s does not exist", configDir)
		scaffolding.Setup(configDir, allowedIPsFile)
	}

	// Load configurations and IP filters
	ipfilter.LoadAllowedIPs(allowedIPsFile)
	config.LoadConfigs(configDir)

	// Start watchers
	go config.WatchConfigDir(configDir)
	go ipfilter.WatchAllowedIPsFile(allowedIPsFile)

	// Setup routes
	app.Get("/*", handler.ConfigHandler)

	// Start server
	applogger.Log.Printf("Starting server on %s", port)
	applogger.Log.Printf("Using config directory: %s", configDir)
	applogger.Log.Printf("Using allowed IPs file: %s", allowedIPsFile)
	applogger.Log.Fatal(app.Listen(port))
}
