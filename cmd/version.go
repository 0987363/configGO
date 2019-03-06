package cmd

import (
	"fmt"

	"github.com/0987363/cobra"
)

// BuildInfo struct is used to contain the information that is set during
// the compilation time
var BuildInfo struct {
	Version string
	Date    string
	Commit  string
	Owner   string
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version information",
	Run:   version,
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("%s %s %s %s\n", BuildInfo.Version, BuildInfo.Date, BuildInfo.Commit, BuildInfo.Owner)
}
