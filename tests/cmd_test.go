//nolint:dupl
package e_tests

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/elcengine/elemental/cmd"
	"github.com/elcengine/elemental/core"
	"github.com/elcengine/elemental/tests/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCmd(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir("..")

	os.Remove(e_cmd.DefaultConfigFile)

	Convey("Call root command", t, func() {
		So(e_cmd.Execute, ShouldNotPanic)
	})

	Convey("Initialize config file", t, func() {
		_, err := os.Stat(e_cmd.DefaultConfigFile)
		So(errors.Is(err, os.ErrNotExist), ShouldBeTrue)

		e_cmd.RootCmd.SetArgs([]string{"init", e_mocks.DEFAULT_DATASOURCE})
		e_cmd.Execute()

		_, err = os.Stat(e_cmd.DefaultConfigFile)
		So(errors.Is(err, os.ErrNotExist), ShouldBeFalse)

		file, err := os.ReadFile(e_cmd.DefaultConfigFile)
		So(err, ShouldBeNil)

		cfg := e_cmd.Config{}
		err = json.Unmarshal(file, &cfg)
		So(err, ShouldBeNil)

		So(cfg.ConnectionStr, ShouldEqual, e_mocks.DEFAULT_DATASOURCE)
	})

	Convey("Migrations and seeds", t, func() {
		checkIfFileExists := func(filename, dir string) bool {
			files, err := os.ReadDir(dir)
			if err != nil {
				return false
			}
			found := false
			for _, file := range files {
				if !file.IsDir() && strings.Contains(file.Name(), filename) {
					found = true
					break
				}
			}
			return found
		}
		Convey("Create migration file", func() {
			filename := "rename_name_to_full_name"
			migrationFileDir := "./database/migrations"

			os.RemoveAll(migrationFileDir)

			So(checkIfFileExists(filename, migrationFileDir), ShouldBeFalse)

			e_cmd.RootCmd.SetArgs([]string{"migrate", "create", filename})
			e_cmd.Execute()

			So(checkIfFileExists(filename, migrationFileDir), ShouldBeTrue)

			Convey("Run migration", func() {
				e_cmd.RootCmd.SetArgs([]string{"migrate", "up"})
				e_cmd.Execute()

				elemental.Connect(e_mocks.DEFAULT_DATASOURCE)

				So(elemental.NativeModel.SetCollection("changelog").
					CountDocuments(primitive.M{"type": "migration"}).ExecInt(), ShouldEqual, 1)

				Convey("Rollback migration", func() {
					e_cmd.RootCmd.SetArgs([]string{"migrate", "down"})
					e_cmd.Execute()

					elemental.Connect(e_mocks.DEFAULT_DATASOURCE)

					So(elemental.NativeModel.SetCollection("changelog").
						CountDocuments(primitive.M{"type": "migration"}).ExecInt(), ShouldEqual, 0)
				})
			})
		})
		Convey("Create seed file", func() {
			filename := "create_test_user"
			seedFileDir := "./database/seeds"

			os.RemoveAll(seedFileDir)

			So(checkIfFileExists(filename, seedFileDir), ShouldBeFalse)

			e_cmd.RootCmd.SetArgs([]string{"seed", "create", filename})
			e_cmd.Execute()

			So(checkIfFileExists(filename, seedFileDir), ShouldBeTrue)

			Convey("Run seed", func() {
				e_cmd.RootCmd.SetArgs([]string{"seed", "up"})
				e_cmd.Execute()

				elemental.Connect(e_mocks.DEFAULT_DATASOURCE)

				So(elemental.NativeModel.SetCollection("changelog").
					CountDocuments(primitive.M{"type": "seed"}).ExecInt(), ShouldEqual, 1)

				Convey("Rollback seed", func() {
					e_cmd.RootCmd.SetArgs([]string{"seed", "down"})
					e_cmd.Execute()

					elemental.Connect(e_mocks.DEFAULT_DATASOURCE)

					So(elemental.NativeModel.SetCollection("changelog").
						CountDocuments(primitive.M{"type": "seed"}).ExecInt(), ShouldEqual, 0)
				})
			})
		})
	})
}
