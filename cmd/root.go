package e_cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "elemental",
	Short: "Your next gen database ODM",
	Long:  `Elemental is a user database ODM that allows you to interact with your database in a much more user friendly way than standard database drivers`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`

Welcome to Elemental!.

------------------------------------		

Please run 'elemental --help' to see available commands.

If you encounter any issues, please report them at "https://github.com/go-elemental/elemental/issues"

------------------------------------`)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(seedCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
