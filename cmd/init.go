package e_cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize elemental with a config file",
	Run: func(cmd *cobra.Command, args []string) {
		configFile := ".elemental.json"
		_, err := os.Stat(configFile)
		if err == nil {
			return
		}
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
		f, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		bytes, _ := json.MarshalIndent(configWithDefaults(&config{}), "", "  ")
		f.Write(bytes)
		fmt.Println("\033[32mElemental config file created at", configFile, "\033[0m")
	},
}
