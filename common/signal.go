package common

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/0987363/configGO/service"
	log "github.com/sirupsen/logrus"
)

func Signal() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
	log.Warning("Recv sig: ", <-sigint)

	service.Exit(0)
}
