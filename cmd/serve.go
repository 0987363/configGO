package cmd

import (
	"github.com/0987363/configGO/handlers"
	"github.com/0987363/configGO/common"
	"github.com/0987363/configGO/middleware"
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

	if cert != "" && key != "" {
		log.Infof("Starting configGO tls server on %s.", address)
		handlers.RootMux.RunTLS(address, cert, key)
	} else {
		log.Infof("Starting configGO server on %s.", address)
		handlers.RootMux.Run(address)
	}
}
