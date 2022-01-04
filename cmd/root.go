package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newCmdRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "nvmet-config",
		Short: "Configure linux nvme target",
	}

	cmdImport := newCmdImport()
	rootCmd.AddCommand(cmdImport)

	return rootCmd
}

// Execute executes comand sequence
func Execute() {
	cmdRoot := newCmdRoot()
	if err := cmdRoot.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
