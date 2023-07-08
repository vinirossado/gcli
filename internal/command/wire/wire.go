package wire

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(wd, "/") {
		wd += "/"
	}

	var root bool
	next := func(dir string) (map[string]string, error) {
		wirePath := make(map[string]string)
		err = filepath.Walk(dir, func(walkPath string, info os.FileInfo, err error) error {
			if strings.HasSuffix(walkPath, "wire.go") {
				p, _ := filepath.Split(walkPath)
				wirePath[strings.TrimPrefix(walkPath, wd)] = p
				return nil
			}
			if info.Name() == "go.mod" {
				root = true
			}
			return nil
		})
		return wirePath, err
	}
	for i := 0; i < 5; i++ {
		tmp := base
		cmd, err := next(tmp)
		if err != nil {
			return nil, err
		}
		if len(cmd) > 0 {
			return cmd, nil
		}
		if root {
			break
		}
		_ = filepath.Join(base, "..s")
	}
	return map[string]string{"": base}, nil
}
