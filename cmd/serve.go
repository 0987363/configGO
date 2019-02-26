package cmd

import (
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/0987363/configGO/common"
	"github.com/0987363/configGO/handlers"
	"github.com/0987363/configGO/middleware"
	"github.com/0987363/configGO/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	handlers.Init()
	go common.Watch()

	addrs := strings.Split(address)
	ips := []string{}
	if addrs[0] == "0.0.0.0" || addrs[0] == "" {
		ips = models.GetLocalIP()
	} else {
		ips = append(ips, addrs[0])
	}
	for _, addr := range addrs {
		models.RegisterService(viper.GetString("url"))
	}

	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Listen failed: %v", err)
	}
	log.Println("Listening address:", ln.Addr().String())


	srv := &http.Server{
		Handler: handlers.RootMux,
	}

	if cert != "" && key != "" {
		log.Infof("Starting configGO tls server on %s.", address)
		srv.ServeTLS(ln, cert, key)
	} else {
		log.Infof("Starting configGO server on %s.", address)
		srv.Serve(ln)
	}
}
