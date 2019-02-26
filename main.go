package main

import (
	"github.com/0987363/configGO/cmd"
	log "github.com/sirupsen/logrus"
	"runtime"
	"time"
)

// Exported onstants for storing build information
var (
	BuildVersion string
	BuildDate    string
	BuildCommit  string
	BuildOwner   string
)

func init() {
	time.Local = time.UTC
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	cmd.BuildInfo.Version = BuildVersion
	cmd.BuildInfo.Date = BuildDate
	cmd.BuildInfo.Commit = BuildCommit
	cmd.BuildInfo.Owner = BuildOwner

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
