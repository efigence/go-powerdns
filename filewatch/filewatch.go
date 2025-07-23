package filewatch

import (
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

type Watcher struct {
}

func NewWatch(log *zap.SugaredLogger, dir string, notifyFunc func(path string)) (err error) {
	dir, err = filepath.Abs(dir)
	if err != nil {
		return err
	}
	watcher, err := fsnotify.NewWatcher()
	//defer watcher.Close()
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Create) {
					st, err := os.Stat(event.Name)
					if err == nil {
						if st.IsDir() && !slices.Contains(watcher.WatchList(), event.Name) {
							log.Infof("adding directory %s watch list", event.Name)
							watcher.Add(event.Name)

						}
					}

				}
				if event.Has(fsnotify.Remove) {
					// I think its automatically removed but still we need message for it
					// so might as well add it just in case
					if slices.Contains(watcher.WatchList(), event.Name) {
						watcher.Remove(event.Name)
					}
					log.Infof("%s removed from watch list", event.Name)
				}
				if event.Has(fsnotify.Write) {
					log.Infof("modified file: %s", event.Name)
					notifyFunc(event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Infof("error: %s", err)
			}
		}
	}()
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			watcher.Add(path)
		}
		return nil
	})
	log.Infof("started watcher on %s", dir)
	return nil
}
