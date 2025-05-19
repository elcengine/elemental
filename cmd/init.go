package cmd

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

		f := lo.Must(os.Create(DefaultConfigFile))
		defer f.Close()

		bytes := lo.Must(json.MarshalIndent(configWithDefaults(&Config{
			ConnectionStr: lo.FirstOrEmpty(args),
		}), "", "  "))
		lo.Must(f.Write(bytes))

		fmt.Println("\033[32mElemental config file created at", DefaultConfigFile, "\033[0m")
	},
}
