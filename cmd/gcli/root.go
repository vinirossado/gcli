package gcli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vinirossado/gcli/internal/command/create"
	"github.com/vinirossado/gcli/internal/command/project"
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
	Version: fmt.Sprintf("%s", version.Version),
}

func init() {
	rootCmd.AddCommand(project.NewCmd)
	rootCmd.AddCommand(create.CreateCmd)
	rootCmd.AddCommand(wire.WireCmd)
	rootCmd.AddCommand(run.Cmd)
	rootCmd.AddCommand(upgrade.UpgradeCmd)
	create.CreateCmd.AddCommand(create.CreateHandlerCmd)
	create.CreateCmd.AddCommand(create.CreateServiceCmd)
	create.CreateCmd.AddCommand(create.CreateRepositoryCmd)
	create.CreateCmd.AddCommand(create.CreateModelCmd)
	create.CreateCmd.AddCommand(create.CreateAllCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
