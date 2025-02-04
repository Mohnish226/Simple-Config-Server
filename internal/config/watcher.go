package config

import (
	"os"
	"path/filepath"
	"simpleConfigServer/internal/logger"

	"github.com/fsnotify/fsnotify"
)

func WatchConfigDir(configDir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Log.Fatal(err)
	}
	defer watcher.Close()

	logger.Log.Printf("Watching config directory: %s", configDir)

	err = filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				logger.Log.Printf("Error adding watcher to directory %s: %v", path, err)
			}
		}
		return nil
	})
	if err != nil {
		logger.Log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				logger.Log.Printf("Config file changed: %s", event.Name)
				LoadConfigFile(event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Log.Printf("Error watching config directory: %v", err)
		}
	}
}
