package common

import (
	"regexp"
	"time"

	"github.com/0987363/configGO/daemon"
	"github.com/0987363/viper"
	"github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"
)

func Watch(c chan *Service) {
	w := NewWorker(viper.GetString("work"), c)
	log.Info("Worker: ", w)
	if w != nil {
		w.ParseProject()
	}

	go func() {
		watch := watcher.New()
		watch.SetMaxEvents(1)

		r := regexp.MustCompile("(?i)json$")
		watch.AddFilterHook(watcher.RegexFilterHook(r, false))

		go func() {
			for {
				select {
				case event := <- watch.Event:
					if event.IsDir() {
						continue
					}

					switch event.Op {
					case watcher.Write:
					case watcher.Create:
						if err := w.Update(event.Path); err != nil {
							log.Error("Update worker failed: ", err)
						}
						log.Warning("file modified: ", event.String())
						continue
					case watcher.Remove:
						if err := w.Remove(event.Path); err != nil {
							log.Error("Update worker failed: ", err)
						}
					default:
						log.Warning("Unknown op: ", event.String())
						continue
					}
				case err := <-watch.Error:
					log.Fatalln(err)
				case <-watch.Closed:
					daemon.Exit()
				}
			}
		}()

		// Watch test_folder recursively for changes.
		if err := watch.AddRecursive(viper.GetString("work")); err != nil {
			log.Fatalln(err)
		}

		if err := watch.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	}()
}
