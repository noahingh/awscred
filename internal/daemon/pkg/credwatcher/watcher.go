package credwatcher

import (
	"context"
	"os"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

// Service watch the event of file.
type Service struct {
	filepath string
}

// NewService create a new watcher.
func NewService(path string) *Service {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Fatalf("the file doesn't exist: %s", path)
	}

	return &Service{
		filepath: path,
	}
}

// Watch send the event when the file has changed until the context is closed.
func (s *Service) Watch(ctx context.Context, ch chan<- fsnotify.Event) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				log.Infof("watcher send the event: %s", event)

				ch <- event

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Errorf("error: %s", err)
			}
		}
	}()

	err = watcher.Add(s.filepath)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("watcher start to watch ...")
	<-ctx.Done()
}
