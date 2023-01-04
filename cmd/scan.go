package cmd

import (
	"github.com/nathanieltruitt/scanner/scanner"
	"github.com/spf13/cobra"
)

var ScanCmd = &cobra.Command{
	Use: "scan",
	Short: "used to scan an IP address or range of IP addresses and output to the specified format.",
	Long: "used to scan an IP address or range of IP addresses and output to the specified format.",
	Run: scan,
}

func scan(cmd *cobra.Command, args []string) {
	scanner.PingRange(Address, Start, End)
}