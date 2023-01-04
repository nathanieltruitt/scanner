package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "moss",
	Short: "Moss is a simple and portable network and port scanner",
	Long: `Moss is a simple and portable network and port scanner written in Go.
				 Written by Nathan Truitt 01/04/23`,
}

var Address string
var Start int
var End int

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	ScanCmd.Flags().StringVarP(&Address, "address", "a", "", "supply an IP or a range of IPs for scanning")
	ScanCmd.Flags().IntVarP(&Start, "start", "s", 1, "starting fourth octet for the scan")
	ScanCmd.Flags().IntVarP(&End, "end", "e", 254, "ending fourth octet for the scan")
	PingCmd.Flags().StringVarP(&Address, "address", "a", "", "ping a single address")
	rootCmd.AddCommand(ScanCmd)
	rootCmd.AddCommand(PingCmd)
}