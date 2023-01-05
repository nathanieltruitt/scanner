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
var Port string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// initialize scan flags
	ScanInit()
	// initialize ping flags
	PingInit()

	// add commands
	rootCmd.AddCommand(ScanCmd)
	rootCmd.AddCommand(PingCmd)
	
	// adding a couple new lines for when the program runs
	fmt.Print("\n\n")
}