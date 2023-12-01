package upgrade

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vinirossado/gcli/config"
	"log"
	"os"
	"os/exec"
)

var CmdUpgrade = &cobra.Command{
	Use:     "upgrade",
	Short:   "Upgrade the gcli command.",
	Long:    "Upgrade the gcli command.",
	Example: "gcli upgrade",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("go install %s\n", config.GcliCmd)
		cmd := exec.Command("go", "install %s\n", config.GcliCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("go install %s error\n", err)
		}
		fmt.Printf("\n🎉, Gcli upgrade successfully!\n\n")
	},
}
