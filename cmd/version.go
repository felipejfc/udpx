package cmd

import (
	"fmt"

	"github.com/felipejfc/udp-proxy/proxy"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of UDP-Proxy",
	Long:  "Print the version number of UDP-Proxy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("UDP-Proxy v%s\n", proxy.VERSION)
	},
}
