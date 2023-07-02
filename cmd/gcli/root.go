package gcli

import (
	"gcli/internal/command/project"
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

}

func Execute() error {
	return rootCmd.Execute()
}
