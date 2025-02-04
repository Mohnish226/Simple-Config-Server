package ipfilter

import (
	"simpleConfigServer/internal/logger"

	"github.com/fsnotify/fsnotify"
)

func WatchAllowedIPsFile(AllowedIPsFile string) {
	logger.Log.Printf("Watching allowed IPs file: %s", AllowedIPsFile)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(AllowedIPsFile)
	if err != nil {
		logger.Log.Fatal("Error watching allowed IPs file: ", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				logger.Log.Printf("Allowed IPs file changed: %s", event.Name)
				LoadAllowedIPs(AllowedIPsFile)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Log.Printf("Error watching allowed IPs file: %v", err)
		}
	}
}
