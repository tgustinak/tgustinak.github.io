package watcher

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

type ProcessFn func() error

const debounceDelay = 300 * time.Millisecond

func Watch(dirs []string, process ProcessFn) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	for _, dir := range dirs {
		if err := addTree(watcher, dir); err != nil {
			return err
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	var timer *time.Timer
	trigger := func(name string) {
		fmt.Printf("File changed: %s\n", name)

		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(debounceDelay, func() {
			start := time.Now()

			if err := process(); err != nil {
				log.Printf("Rebuild failed: %v\n", err)

				return
			}

			fmt.Printf("Rebuilt in %s\n", time.Since(start).Round(time.Millisecond))
		})
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename|fsnotify.Chmod) == 0 {
				continue
			}

			if event.Op&fsnotify.Create != 0 {
				if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
					if err := addTree(watcher, event.Name); err != nil {
						log.Printf("Watcher error: %v\n", err)
					}
				}
			}

			trigger(event.Name)
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}

			log.Printf("Watcher error: %v\n", err)
		case <-signals:
			fmt.Println("\nStopping watcher.")

			return nil
		}
	}
}

func addTree(watcher *fsnotify.Watcher, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		return watcher.Add(path)
	})
}
