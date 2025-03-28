package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type AuditEvent struct {
	Timestamp   string                 `json:"timestamp"`
	EventType   string                 `json:"event_type"`
	ClientIP    string                 `json:"client_ip"`
	UserID      string                 `json:"user_id,omitempty"`
	Product     string                 `json:"product,omitempty"`
	Environment string                 `json:"environment,omitempty"`
	ConfigKey   string                 `json:"config_key,omitempty"`
	Status      string                 `json:"status"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

var auditFile *os.File

func init() {
	// Create audit directory if it doesn't exist
	auditDir := "audit_logs"
	if err := os.MkdirAll(auditDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create audit directory: %v", err))
	}

	// Create new audit log file with date
	filename := filepath.Join(auditDir, fmt.Sprintf("audit_%s.log", time.Now().Format("2006-01-02")))
	var err error
	auditFile, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open audit log file: %v", err))
	}
}

func Log(eventType string, clientIP string, status string, details map[string]interface{}) {
	event := AuditEvent{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		EventType: eventType,
		ClientIP:  clientIP,
		Status:    status,
		Details:   details,
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Failed to marshal audit event: %v\n", err)
		return
	}

	if _, err := auditFile.Write(append(jsonData, '\n')); err != nil {
		fmt.Printf("Failed to write audit log: %v\n", err)
	}
}

func LogAuth(clientIP string, status string, userID string) {
	details := map[string]interface{}{
		"user_id": userID,
	}
	Log("AUTH", clientIP, status, details)
}

func LogConfigAccess(clientIP string, status string, product string, env string, configKey string, userID string) {
	details := map[string]interface{}{
		"product":     product,
		"environment": env,
		"config_key":  configKey,
		"user_id":     userID,
	}
	Log("CONFIG_ACCESS", clientIP, status, details)
}

func LogConfigChange(clientIP string, status string, product string, env string, configKey string, oldValue string, newValue string, userID string) {
	details := map[string]interface{}{
		"product":     product,
		"environment": env,
		"config_key":  configKey,
		"old_value":   oldValue,
		"new_value":   newValue,
		"user_id":     userID,
	}
	Log("CONFIG_CHANGE", clientIP, status, details)
}

func LogSecurity(clientIP string, status string, eventType string, details map[string]interface{}) {
	Log(eventType, clientIP, status, details)
}

func LogSystem(eventType string, status string, details map[string]interface{}) {
	Log(eventType, "SYSTEM", status, details)
}
