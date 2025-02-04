package config

import (
	"os"
	"path/filepath"
	"simpleConfigServer/internal/logger"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

var configStore = make(map[string]map[string]string)
var configLoadMux sync.Mutex
var mu sync.RWMutex

type Config struct {
	Configs map[string]string `yaml:"configs"`
}

func LoadConfigs(configPath string) {

	configLoadMux.Lock()
	defer configLoadMux.Unlock()

	err := filepath.Walk(configPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".yml" {
			LoadConfigFile(path)
			logger.Log.Printf("Loaded config file: %s", path)
		}
		return nil
	})

	if err != nil {
		logger.Log.Fatalf("Error walking config directory: %v", err)
	}
}

func LoadConfigFile(path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		logger.Log.Printf("Failed to read %s: %v", path, err)
		return
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		logger.Log.Printf("Failed to parse YAML %s: %v", path, err)
		return
	}

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		logger.Log.Printf("Invalid path structure: %s", path)
		return
	}

	product := filepath.Base(filepath.Dir(path))
	mu.Lock()
	configStore[product] = config.Configs
	mu.Unlock()
	logger.Log.Printf("Loaded configs for %s", product)
}

func GetConfigs() map[string]map[string]string {
	return configStore
}
