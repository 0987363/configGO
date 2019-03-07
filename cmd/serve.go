package cmd

import (
	"github.com/0987363/configGO/common"
	"github.com/0987363/configGO/middleware"
	"github.com/0987363/configGO/daemon"
	"github.com/spf13/cobra"
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

	//	handlers.Init()

	c := common.Registry()
	common.Watch(c)
	common.Signal()

	//	common.Registry(ln.Addr().String())

	daemon.Wait()
}
