package main

import (
	"net/http"
	"os"
	"simpleConfigServer/internal/config"
	"simpleConfigServer/internal/handler"
	"simpleConfigServer/internal/ipfilter"
	"simpleConfigServer/internal/logger"
	"simpleConfigServer/internal/scaffolding"
)

var ConfigDir = func() string {
	if dir := os.Getenv("CONFIG_DIR"); dir != "" {
		return dir
	}
	return "configurations"
}()

var AllowedIPsFile = func() string {
	if file := os.Getenv("ALLOWED_IPS_FILE"); file != "" {
		return file
	}
	return "allowed_ips.txt"
}()

var port = func() string {
	if p := os.Getenv("PORT"); p != "" {
		return ":" + p
	}
	return ":8080"
}()

func main() {

	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		logger.Log.Printf("Configurations directory %s does not exist", ConfigDir)
		scaffolding.Setup(ConfigDir, AllowedIPsFile)
	}

	ipfilter.LoadAllowedIPs(AllowedIPsFile)
	config.LoadConfigs(ConfigDir)

	go config.WatchConfigDir(ConfigDir)
	go ipfilter.WatchAllowedIPsFile(AllowedIPsFile)

	http.HandleFunc("/", handler.ConfigHandler)
	logger.Log.Printf("Starting server on %s", port)
	logger.Log.Fatal(http.ListenAndServe(port, nil))
}
