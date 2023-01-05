package cmd

import (
	"fmt"

	"github.com/nathanieltruitt/scanner/scanner"
	"github.com/spf13/cobra"
)


var PingCmd = &cobra.Command{
	Use: "ping",
	Short: "ping a single address",
	Long: "Ping a single address",
	Run: ping,
}

func PingInit() {
	PingCmd.Flags().StringVarP(&Address, "address", "a", "", "ping a single address")
}

func ping(cmd *cobra.Command, args []string) {
	data := scanner.ScanData{}
	scanner.PingSingle(Address, &data)
	if len(data.OnlineHosts) > 0 {
		fmt.Println("ping successful: ", data.OnlineHosts[0])
	}
}