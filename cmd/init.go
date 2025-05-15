package e_cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

const DefaultConfigFile = ".elementalrc"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize elemental with a config file",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat(DefaultConfigFile)
		if err == nil {
			return
		}
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
		f, err := os.Create(DefaultConfigFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		bytes, err := json.MarshalIndent(configWithDefaults(&Config{
			ConnectionStr: lo.FirstOrEmpty(args),
		}), "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\033[32mElemental config file created at", DefaultConfigFile, "\033[0m")
	},
}
