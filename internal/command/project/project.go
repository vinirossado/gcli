package project

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

type Project struct {
	ProjectName string `survey:"name"`
}

var NewCmd = &cobra.Command{
	Use:     "new",
	Example: "gcli new awesome-api",
	Short:   "Create a new Project Api",
	Long:    "Create a new project with GCLI layout",
	Run:     run,
}

var (
	repoURL string
)

func init() {
	NewCmd.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
}

func run(cmd *cobra.Command, args []string) {
	p := NewProject()
	if len(args) == 0 {
		err := survey.AskOne(&survey.Input{
			Message: "What is your project name",
			Help:    "project name.",
			Suggest: nil,
		}, &p.ProjectName, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println(err.Error())
			return
		} else {
			p.ProjectName = args[0]
		}
	}
}

func NewProject() *Project {
	return &Project{}
}
