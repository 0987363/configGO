package cmd

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/0987363/configGO/common"
	"github.com/0987363/configGO/handlers"
	"github.com/0987363/configGO/middleware"
	"github.com/0987363/configGO/service"
	log "github.com/sirupsen/logrus"
	"github.com/0987363/cobra"
	"github.com/0987363/viper"
)

const defaultAddress = ":10001"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:    "serve",
	Short:  "Start config server",
	PreRun: LoadConfiguration,
	Run:    serve,
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	middleware.LoggerConnInit()

	address := viper.GetString("address")
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")

	addrs := strings.Split(address, ":")
	if len(addrs) != 2 {
		log.Fatal("Address failed:", address)
	}

	handlers.Init()

	common.Watch()
	common.Signal()

	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Listen failed:", err)
	}
	log.Println("Listening address:", ln.Addr().String())
	common.Registry(ln.Addr().String())

	srv := &http.Server{
		Handler: handlers.RootMux,
	}
	idleConnsClosed := make(chan struct{})
	service.AddCloseHook(func() {
		srv.Shutdown(context.Background())
		close(idleConnsClosed)
	})
	if cert != "" && key != "" {
		log.Infof("Starting configGO tls server on %s.", address)
		srv.ServeTLS(ln, cert, key)
	} else {
		log.Infof("Starting configGO server on %s.", address)
		srv.Serve(ln)
	}
	<-idleConnsClosed
}
