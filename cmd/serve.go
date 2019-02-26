package cmd

import (
	"github.com/0987363/configGO/handlers"
	"github.com/0987363/configGO/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultAddress = ":8080"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:    "serve",
	Short:  "Start cattle server",
	PreRun: LoadConfiguration,
	Run:    serve,
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	middleware.LoggerConnInit()
	/*
	middleware.ConnectMicro()

	// Try to connect to the database
	if err := middleware.ConnectDB(
		viper.GetString("database.mongodb"),
	); err != nil {
		log.Fatalln("connect to db:", err)
	}
	*/

	address := viper.GetString("address")
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")

	handlers.Init()

	if cert != "" && key != "" {
		log.Infof("Starting cattle admin tls server on %s.", address)
		handlers.RootMux.RunTLS(address, cert, key)
	} else {
		log.Infof("Starting cattle admin server on %s.", address)
		handlers.RootMux.Run(address)
	}
}
