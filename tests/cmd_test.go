package e_tests

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	e_cmd "github.com/elcengine/elemental/cmd"
	e_mocks "github.com/elcengine/elemental/tests/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCmd(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir("..")

	os.Remove(e_cmd.DefaultConfigFile)

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
}
