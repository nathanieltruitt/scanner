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

func ping(cmd *cobra.Command, args []string) {
	data := scanner.ScanData{}
	scanner.Ping(Address, &data)
	fmt.Println(data.OnlineHosts)
}