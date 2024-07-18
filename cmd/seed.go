package e_cmd

import (
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Run database seeds",
}

var createSeedCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new seed",
	Run: func(cmd *cobra.Command, args []string) {
		create(args, "seed")
	},
}

var runSeedCmd = &cobra.Command{
	Use:   "up",
	Short: "Run database seeds",
	Run: func(cmd *cobra.Command, args []string) {
		run(false, "seed")
	},
}

var rollbackSeedCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback database seeds",
	Run: func(cmd *cobra.Command, args []string) {
		run(true, "seed")
	},
}

func init() {
	seedCmd.AddCommand(createSeedCmd)
	seedCmd.AddCommand(runSeedCmd)
	seedCmd.AddCommand(rollbackSeedCmd)
}
