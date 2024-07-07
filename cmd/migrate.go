package e_cmd

import (
	e_connection "elemental/connection"
	"fmt"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := readConfigFile()
		e_connection.ConnectURI(cfg.ConnectionStr)
		defer e_connection.Disconnect()
		fmt.Println("Running database migrations", cfg)
	},
}
