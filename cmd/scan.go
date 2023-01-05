package cmd

import (
	"fmt"

	"github.com/nathanieltruitt/scanner/scanner"
	"github.com/spf13/cobra"
)

var ScanCmd = &cobra.Command{
	Use: "scan",
	Short: "used to scan an IP address or range of IP addresses and output to the specified format.",
	Long: "used to scan an IP address or range of IP addresses and output to the specified format.",
	Run: scan,
}

func ScanInit() {
	ScanCmd.Flags().StringVarP(&Address, "address", "a", "", "supply an IP or a range of IPs for scanning")
	ScanCmd.Flags().StringVarP(&Port, "port", "p", "1-1000", "supply a port to scan")
	ScanCmd.Flags().IntVarP(&Start, "start", "s", 1, "starting fourth octet for the scan")
	ScanCmd.Flags().IntVarP(&End, "end", "e", 254, "ending fourth octet for the scan")
}

func scan(cmd *cobra.Command, args []string) {
	data := scanner.ScanData{}
	scanner.PingRange(Address, Start, End, &data)
	for _, addr := range data.OnlineHosts {
		scanner.Scan(addr, Port, "tcp")
	}
	// add a new line after output
	fmt.Print("\n")
}