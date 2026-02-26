package gcli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vinirossado/gcli/internal/command/create"
	"github.com/vinirossado/gcli/internal/command/new"
	"github.com/vinirossado/gcli/internal/command/run"
	"github.com/vinirossado/gcli/internal/command/upgrade"
	"github.com/vinirossado/gcli/internal/command/wire"
	"github.com/vinirossado/gcli/internal/version"
)

var art = ".----------------.  .----------------.  .----------------.  .----------------.\n" +
	"| .--------------. || .--------------. || .--------------. || .--------------. |\n" +
	"| |    ______    | || |     ______   | || |   _____      | || |     _____    | |\n" +
	"| |  .' ___  |   | || |   .' ___  |  | || |  |_   _|     | || |    |_   _|   | |\n" +
	"| | / .'   \\_|   | || |  / .'   \\_|  | || |    | |       | || |      | |     | |\n" +
	"| | | |    ____  | || |  | |         | || |    | |   _   | || |      | |     | |\n" +
	"| | \\ `.___]  _| | || |  \\ `.___.'\\  | || |   _| |__/ |  | || |     _| |_    | |\n" +
	"| |  `._____.'   | || |   `._____.'  | || |  |________|  | || |    |_____|   | |\n" +
	"| |              | || |              | || |              | || |              | |\n" +
	"| '--------------' || '--------------' || '--------------' || '--------------' |\n" +
	"'----------------'  '----------------'  '----------------'  '----------------'"

var rootCmd = &cobra.Command{
	Use:     "gcli",
	Example: "gcli new awesome-api",
	Short:   fmt.Sprintf("\n%s", art),
	Version: version.Version,
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(new.NewCmd)
	rootCmd.AddCommand(create.CreateCmd)
	rootCmd.AddCommand(run.CmdRun)
	rootCmd.AddCommand(upgrade.CmdUpgrade)

	create.CreateCmd.AddCommand(create.CmdCreateHandler)
	create.CreateCmd.AddCommand(create.CmdCreateService)
	create.CreateCmd.AddCommand(create.CmdCreateRepository)
	create.CreateCmd.AddCommand(create.CmdCreateModel)
	create.CreateCmd.AddCommand(create.CmdCreateAll)

	rootCmd.AddCommand(wire.CmdWire)
	wire.CmdWire.AddCommand(wire.CmdWireAll)
}

func Execute() error {
	return rootCmd.Execute()
}
