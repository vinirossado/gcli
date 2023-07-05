package run

import (
	"fmt"
	"gcli/internal/pkg/helper"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var quit = make(chan os.Signal, 1)

type Run struct{}

var RunCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run",
	Long:    "Run",
	Example: "FindMain",
	Run: func(cmd *cobra.Command, args []string) {
		cmdArgs, programArgs := helper.SplitArgs(cmd, args)
		var dir string
		if len(cmdArgs) > 0 {
			dir = cmdArgs[0]
		}
		base, err := os.Getwd()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
			return
		}
		if dir == "" {
			cmdPath, err := helper.FindMain(base)

			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
				return
			}
			switch len(cmdPath) {
			case 0:
				_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", "The cmd directory cannot be found in the current working directory")
				return
			case 1:
				for _, v := range cmdPath {
					dir = v
				}
			default:
				var cmdPaths []string
				for k := range cmdPath {
					cmdPaths = append(cmdPaths, k)
				}
				prompt := &survey.Select{
					Message:  "Which directory do you want to run?",
					Options:  cmdPaths,
					PageSize: 10,
				}
				e := survey.AskOne(prompt, &dir)
				if e != nil || dir == "" {
					return
				}
				dir = cmdPath[dir]
			}
		}
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		fmt.Printf("\033[35mGcli run %s.\033[0m\n", dir)
		watch(dir, programArgs)
	},
}

func watch(dir string, programArgs []string) {
	print("Watching")
	print(dir)
	for i := 0; i < len(programArgs); i++ {
		print(programArgs[i])
	}
}
