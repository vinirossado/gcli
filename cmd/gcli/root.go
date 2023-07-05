package gcli

import (
	"gcli/internal/command/create"
	"gcli/internal/command/project"
	"gcli/internal/command/run"
	"gcli/internal/command/upgrade"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gcli",
	Example: "gcli new awesome-app",
	Short:   "\n Create a new awesome",
	Version: "AAAAA",
}

func init() {
	rootCmd.AddCommand(project.NewCmd)
	rootCmd.AddCommand(create.CreateCmd)
	//rootCmd.AddCommand(wire.)
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
