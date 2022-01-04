package cmd

import (
	"log"

	"github.com/0xfd4d/nvmet-config/nvmet"

	"github.com/spf13/cobra"
)

func newCmdImport() *cobra.Command {
	var importCmd = &cobra.Command{
		Use:   "import <file>",
		Short: "Apply file state to target",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			nvmf := nvmet.Nvmf{}
			if err := nvmf.ReadFile(args[0]); err != nil {
				log.Fatalln(err)
			}
			if err := nvmf.Apply(); err != nil {
				log.Fatalln(err)
			}
		},
	}

	return importCmd
}
