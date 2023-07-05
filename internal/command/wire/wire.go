package wire

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
)

var WireCmd = &cobra.Command{
	Use:     "wire",
	Short:   "gcli wire [wire.go path]",
	Long:    "gcli wire [wire.go path]",
	Example: "gcli wire cmd/server",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func wire(wirePath string) {
	fmt.Println("wire.go path: ", wirePath)
	cmd := exec.Command("wire")
	cmd.Dir = wirePath
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("wire fail", err)
	}
	fmt.Println(string(out))
}

func findWire(base string) (map[string]string, error) {
	return map[string]string{"": base}, nil
}
