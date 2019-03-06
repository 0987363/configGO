package common

import (
	"regexp"
	"syscall"
	"time"

	"github.com/0987363/configGO/service"
	"github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"
	"github.com/0987363/viper"
)

func Watch() {
	go func() {
		w := watcher.New()
		w.SetMaxEvents(1)

		r := regexp.MustCompile("(?i)(json$)|(toml$)")
		w.AddFilterHook(watcher.RegexFilterHook(r, false))

		go func() {
			for {
				select {
				case event := <-w.Event:
					switch event.Op {
					case watcher.Write:
					case watcher.Create:
//					case watcher.Remove:
//					case watcher.Rename:
//					case watcher.Move:

					}
					syscall.Kill(syscall.Getpid(), syscall.SIGUSR2)
//					syscall.Tgkill(syscall.Getpid(), syscall.Gettid(), syscall.SIGUSR2)
					log.Warning("Recv file changed event: ", event)
					continue
					service.Exit(0)
				case err := <-w.Error:
					log.Fatalln(err)
				case <-w.Closed:
					return
				}
			}
		}()

		// Watch test_folder recursively for changes.
		if err := w.AddRecursive(viper.GetString("work")); err != nil {
			log.Fatalln(err)
		}

		if err := w.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	}()
}
