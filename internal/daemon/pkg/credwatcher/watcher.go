package credwatcher

import (
	"context"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
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

				log.Println("watcher send the event:", event)
				ch <- event

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(s.filepath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("watcher start to watch ...")
	<-ctx.Done()
}
