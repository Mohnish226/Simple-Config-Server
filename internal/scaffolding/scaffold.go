package scaffolding

import (
	"fmt"
	"os"
	"path/filepath"
	"simpleConfigServer/internal/logger"
)

const defaultYAML = `configs:
  version: 1.0
  feature_1: true
  feature_2: false
  logging_level: debug
  logging_file: development.log
  snmp_host: "localhost"
  snmp_port: 161
  snmp_community: "public"
  snmp_version: "v2c"
  snmp_timeout: 1
  snmp_retries: 5
  snmp_max_repetitions: 25
  snmp_non_repeaters: 0
`

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func ensureFile(filePath, content string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return os.WriteFile(filePath, []byte(content), 0644)
	}
	return nil
}

func Setup(baseDir string, allowedIPsFile string) {
	_, err := os.Stat(baseDir)

	if os.IsNotExist(err) {
		logger.Log.Printf("Configurations directory %s does not exist", baseDir)

		if err := ensureDir(baseDir); err != nil {
			logger.Log.Fatalf("Failed to create directory %s: %v", baseDir, err)
		}

		var sampleDir = filepath.Join(baseDir, "sample")

		if err := ensureDir(sampleDir); err != nil {
			logger.Log.Fatalf("Failed to create directory %s: %v", sampleDir, err)
		}

		var configFile = filepath.Join(sampleDir, "development.yml")

		if err := ensureFile(configFile, defaultYAML); err != nil {
			fmt.Printf("Failed to create file %s: %v\n", configFile, err)
		}
	} else {
		logger.Log.Printf("Configurations directory %s already exists", baseDir)
	}

	// Check if allowed_ips.txt exists on the root directory
	_, err = os.Stat(allowedIPsFile)
	if os.IsNotExist(err) {
		logger.Log.Printf("Allowed IPs file %s does not exist", allowedIPsFile)

		if err := ensureFile(allowedIPsFile, ""); err != nil {
			logger.Log.Fatalf("Failed to create file %s: %v", allowedIPsFile, err)
		}
	} else {
		logger.Log.Printf("Allowed IPs file %s already exists", allowedIPsFile)
	}

}
