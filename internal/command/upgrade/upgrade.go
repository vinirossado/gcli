package upgrade

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/vinirossado/gcli/config"
)

var CmdUpgrade = &cobra.Command{
	Use:     "upgrade",
	Short:   "Upgrade gcli to the latest version.",
	Long:    "Upgrade gcli to the latest version.",
	Example: "gcli upgrade",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("go install %s\n", config.GcliCmd)
		cmd := exec.Command("go", "install", config.GcliCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("gcli upgrade failed: %v", err)
		}
		fmt.Println("\ngcli upgraded successfully!")
	},
}
