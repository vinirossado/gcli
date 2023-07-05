package gcli

import (
	"fmt"
	"gcli/internal/command/create"
	"gcli/internal/command/project"
	"gcli/internal/command/run"
	"gcli/internal/command/upgrade"
	"gcli/internal/command/wire"
	"gcli/internal/version"
	"github.com/spf13/cobra"
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
	rootCmd.AddCommand(run.RunCmd)
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
