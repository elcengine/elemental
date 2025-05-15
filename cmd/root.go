package e_cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"maps"
)

type Config struct {
	ConnectionStr       string `json:"connection_str" yaml:"connection_str"`             // The connection string to connect with the data source
	MigrationsDir       string `json:"migrations_dir" yaml:"migrations_dir"`             // The directory where migration files are stored
	SeedsDir            string `json:"seeds_dir" yaml:"seeds_dir"`                       // The directory where seed files are stored
	ChangelogCollection string `json:"changelog_collection" yaml:"changelog_collection"` // The collection where changelogs are stored in the database
}

var RootCmd = &cobra.Command{
	Use:   "elemental",
	Short: "Your next gen MongoDB ODM",
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
	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(migrateCmd)
	RootCmd.AddCommand(seedCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func configWithDefaults(conf *Config) Config {
	if conf.MigrationsDir == "" {
		conf.MigrationsDir = "database/migrations"
	}
	if conf.SeedsDir == "" {
		conf.SeedsDir = "database/seeds"
	}
	if conf.ChangelogCollection == "" {
		conf.ChangelogCollection = "changelog"
	}
	return *conf
}

func readConfigFile() Config {
	var conf Config

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	supportedConfigFiles := []string{".elementalrc", "elemental.json", "elemental.yaml", "elemental.yml", ".elemental.yaml", ".elemental.yml"}

	for i, file := range supportedConfigFiles {
		configFilePath := filepath.Join(dir, file)
		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			if i == len(supportedConfigFiles)-1 {
				log.Fatalf(`Config file not found. Please create a config file matching one of the following names: %s 
				and place it in the root of your project or run 'elemental init' to create one.`, strings.Join(supportedConfigFiles, ", "))
			}
			continue
		}
		file, err := os.ReadFile(configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasSuffix(configFilePath, ".yaml") || strings.HasSuffix(configFilePath, ".yml") {
			err = yaml.Unmarshal(file, &conf)
		} else {
			err = json.Unmarshal(file, &conf)
		}
		if err != nil {
			log.Fatal("Failed to decode config file with error:", err)
		}
		break
	}
	if conf.ConnectionStr == "" {
		log.Fatal("Connection string is required in the config file")
	}

	envFile := ".env"

	appEnv := os.Getenv("APP_ENV")
	if appEnv != "" {
		envFile = fmt.Sprintf(".env.%s", appEnv)
	}
	envVars := map[string]string{}

	if _, err := os.Stat(envFile); err == nil {
		fileVars, _ := godotenv.Read(envFile)
		maps.Copy(envVars, fileVars)
	}

	if strings.HasPrefix(conf.ConnectionStr, "<") && strings.HasSuffix(conf.ConnectionStr, ">") {
		key := strings.Trim(conf.ConnectionStr, "<>")
		conf.ConnectionStr = lo.CoalesceOrEmpty(envVars[key], os.Getenv(key))
	}

	return configWithDefaults(&conf)
}
