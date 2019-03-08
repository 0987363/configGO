package common

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/0987363/configGO/daemon"
	log "github.com/sirupsen/logrus"
)

func Signal() {
	go func() {
		for {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
			switch <-sigint {
			case syscall.SIGUSR2:
				log.Warning("Recv restart sig.")
				daemon.Exit()
				//	service.Restart()
				continue
			default:
				log.Warning("Recv close sig.")
				daemon.Exit()
			}
		}
	}()
}
