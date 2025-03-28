package ipfilter

import (
	"simpleConfigServer/internal/audit"
	"simpleConfigServer/internal/logger"

	"github.com/fsnotify/fsnotify"
)

func WatchAllowedIPsFile(AllowedIPsFile string) {
	logger.Log.Printf("Watching allowed IPs file: %s", AllowedIPsFile)
	audit.LogSystem("IP_FILTER_WATCH", "STARTED", map[string]interface{}{
		"file": AllowedIPsFile,
	})

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Log.Fatal(err)
		audit.LogSystem("IP_FILTER_WATCH", "FAILED", map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer watcher.Close()

	err = watcher.Add(AllowedIPsFile)
	if err != nil {
		logger.Log.Fatal("Error watching allowed IPs file: ", err)
		audit.LogSystem("IP_FILTER_WATCH", "FAILED", map[string]interface{}{
			"error": err.Error(),
		})
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				logger.Log.Printf("Allowed IPs file changed: %s", event.Name)
				audit.LogSystem("IP_FILTER_CHANGE", "DETECTED", map[string]interface{}{
					"file": event.Name,
					"op":   event.Op.String(),
				})
				LoadAllowedIPs(AllowedIPsFile)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Log.Printf("Error watching allowed IPs file: %v", err)
			audit.LogSystem("IP_FILTER_WATCH", "ERROR", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
}
