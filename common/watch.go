package common

import (
	"regexp"
	"time"

	"github.com/0987363/configGO/daemon"
	"github.com/0987363/viper"
	"github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"
)

func Watch(c chan *Application) {
	w := NewWorker(viper.GetString("work"))
	log.Info("Read first worker:", w)

	if w != nil {
		w.ParseProject(c)
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
					switch event.Op {
					case watcher.Write:
						if err := w.Update(event.Path, c); err != nil {
							log.Error("Update worker failed: ", err)
						}
						log.Warning("file modified: ", event.Name(), event.Path, event.String(), event.FileInfo)
						continue
					case watcher.Create:
						log.Warning("file create: ", event.Name(), event.Path, event.String(), event.FileInfo)
						continue
						//					case watcher.Remove:
					default:
						continue
					}
					//		syscall.Kill(syscall.Getpid(), syscall.SIGUSR2)
					//					syscall.Tgkill(syscall.Getpid(), syscall.Gettid(), syscall.SIGUSR2)
					//		log.Warning("Recv file changed event: ", event)
					//		continue
					//		service.Exit()
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
