//nolint:gocritic
package e_cmd

import (
	"fmt"
	"github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/utils"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

const (
	TargetSeed    = "seed"
	TargetMigrate = "migrate"
)

func run(rollback bool, target string) {
	cfg := readConfigFile()
	dir := cfg.MigrationsDir
	if target == TargetSeed {
		dir = cfg.SeedsDir
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read %ss: %v", target, err)
	}
	extractTimestamp := func(fileName string) int64 {
		fileName = strings.TrimSuffix(fileName, ".go")
		parts := strings.Split(fileName, "_")
		timestampStr := parts[len(parts)-1]
		return lo.Must(strconv.ParseInt(timestampStr, 10, 64))
	}
	sort.Slice(files, func(i, j int) bool {
		file1Timestamp := extractTimestamp(files[i].Name())
		file2Timestamp := extractTimestamp(files[j].Name())
		return file1Timestamp < file2Timestamp
	})

	var template = fmt.Sprintf(`package main 

import (
	"context"
	"strings"
	"strconv"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/samber/lo"
	"github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/database/%ss"
)

 func extractTimestamp(fileName string) int64 {
		fileName = strings.TrimSuffix(fileName, ".go")
		parts := strings.Split(fileName, "_")
		timestampStr := parts[len(parts)-1]
		return lo.Must(strconv.ParseInt(timestampStr, 10, 64))
}

func main() {
	client := elemental.Connect("%s")
	db := elemental.UseDefaultDatabase()
	go db.Collection("%s").Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: map[string]interface{}{"type": 1},
	})
	files := []string{%s}
	functions := []func(context.Context, *mongo.Database, *mongo.Client){%s}
	ctx := context.Background()
	for i, f := range files {
		functionToRun := functions[i]
		var entry = map[string]interface{}{}
		db.Collection("%s").FindOne(ctx, map[string]interface{}{"name": f, "type": "%s"}).Decode(&entry)
		if %t {
			if entry["name"] != nil {
				functionToRun(ctx, db, &client)
				db.Collection("%s").DeleteOne(ctx, map[string]interface{}{"name": f})
			}
		} else {
			if entry["name"] == nil {
				functionToRun(ctx, db, &client)
				db.Collection("%s").InsertOne(ctx, map[string]interface{}{
					"name": f,
					"type": "%s",
					"created_at": time.Now(),
				})
			} 
		}
	}
}
`,
		target, cfg.ConnectionStr, cfg.ChangelogCollection,
		strings.Join(lo.Map(files, func(file os.DirEntry, index int) string {
			return fmt.Sprintf("\"%s\"", strings.TrimSuffix(file.Name(), ".go"))
		}), ","),
		strings.Join(lo.Map(files, func(file os.DirEntry, index int) string {
			if rollback {
				return fmt.Sprintf("%ss.Down_%d", target, extractTimestamp(file.Name()))
			} else {
				return fmt.Sprintf("%ss.Up_%d", target, extractTimestamp(file.Name()))
			}
		}), ","),
		cfg.ChangelogCollection,
		target,
		rollback, cfg.ChangelogCollection, cfg.ChangelogCollection, target,
	)
	elemental.Connect(cfg.ConnectionStr)
	defer elemental.Disconnect()
	os.MkdirAll(".elemental/"+target+"s", os.ModePerm)
	e_utils.CreateAndWriteToFile(fmt.Sprintf(".elemental/%ss/main.go", target), template)
	err = exec.Command("go", "run", ".elemental/"+target+"s/main.go").Run()
	if err != nil {
		elemental.Disconnect()
		log.Fatalf("Failed to run %ss: %s", target, err.Error())
	} else {
		if rollback {
			log.Printf("Successfully rolled back %ss", target)
		} else {
			log.Printf("Successfully ran %ss", target)
		}
	}
	os.Remove(fmt.Sprintf(".elemental/%ss/main.go", target))
}

func create(args []string, target string) {
	if len(args) == 0 {
		log.Fatalf("Please provide a name for the %s", target)
	}
	cfg := readConfigFile()
	if target != TargetSeed {
		os.MkdirAll(cfg.MigrationsDir, os.ModePerm)
	} else {
		os.MkdirAll(cfg.SeedsDir, os.ModePerm)
	}
	timestamp := cast.ToString(time.Now().UnixMilli())
	var template = fmt.Sprintf(`package %ss

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

func Up_%s(ctx context.Context, db *mongo.Database, client *mongo.Client) {
	// Write your %s here
}

func Down_%s(ctx context.Context, db *mongo.Database, client *mongo.Client) {
	// Write your rollback here
}`, target, timestamp, target, timestamp)
	dir := cfg.MigrationsDir
	if target == TargetSeed {
		dir = cfg.SeedsDir
	}
	outputFile := dir + "/" + args[0] + "_" + timestamp + ".go"
	e_utils.CreateAndWriteToFile(outputFile, template)
	log.Printf("%s file created at %s", target, outputFile)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
}

var createMigrationCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new migration",
	Run: func(cmd *cobra.Command, args []string) {
		create(args, "migration")
	},
}

var runMigrationCmd = &cobra.Command{
	Use:   "up",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		run(false, "migration")
	},
}

var rollbackMigrationCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		run(true, "migration")
	},
}

func init() {
	migrateCmd.AddCommand(createMigrationCmd)
	migrateCmd.AddCommand(runMigrationCmd)
	migrateCmd.AddCommand(rollbackMigrationCmd)
}
