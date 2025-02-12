package new

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/vinirossado/gcli/config"
	"github.com/vinirossado/gcli/internal/pkg/helper"
)

type Project struct {
	ProjectName string `survey:"name"`
}

var NewCmd = &cobra.Command{
	Use:     "new",
	Example: "gcli new awesome-api",
	Short:   "Create a new Project Api",
	Long:    "Create a new new with GCLI layout",
	Run:     run,
}

var (
	repoURL string
)

func init() {
	NewCmd.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
}

func NewProject() *Project {
	return &Project{}
}

func run(cmd *cobra.Command, args []string) {
	p := NewProject()
	if len(args) == 0 {
		err := survey.AskOne(&survey.Input{
			Message: "What is your new name?",
			Help:    "new name.",
			Suggest: nil,
		}, &p.ProjectName, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		p.ProjectName = args[0]
	}

	yes, err := p.cloneTemplate()
	if err != nil || !yes {
		return
	}

	err = p.replacePackageName()
	if err != nil || !yes {
		return
	}

	err = p.replacePackageName()
	if err != nil || !yes {
		return
	}

	err = p.modTidy()
	if err != nil || !yes {
		return
	}

	p.rmGit()
	p.installWire()
	fmt.Printf("\n\n🎉 Project \u001B[36m%s\u001B[0m created successfully!\n", p.ProjectName)
	fmt.Printf("\n🎉 Setup DB and run and Docker compose\n\n")
	fmt.Printf("Done. Now run:\n\n")
	fmt.Printf("› \033[36mcd %s \033[0m\n", p.ProjectName)
	fmt.Printf("\n› \033[36mgcli run \033[0m\n")
}

func (p Project) installWire() {
	fmt.Printf("go install %s\n", config.WireCmd)
	cmd := exec.Command("go", "install", config.WireCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("go install %s error \n", err)
	}
}

func (p *Project) modTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = p.ProjectName
	if err := cmd.Run(); err != nil {
		fmt.Println("go mod tidy error: ", err)
		return err
	}
	return nil
}

func (p *Project) rmGit() {
	os.RemoveAll(p.ProjectName + "/.git")
}

func (p *Project) cloneTemplate() (bool, error) {
	stat, _ := os.Stat(p.ProjectName)
	if stat != nil {
		var overwrite = false

		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Folder %s alrerady exists, do you want to overwrite it?", p.ProjectName),
			Help:    "Remove old new and create new new",
		}

		err := survey.AskOne(prompt, &overwrite)
		if err != nil {
			return false, err
		}
		if !overwrite {
			return false, nil
		}
		err = os.RemoveAll(p.ProjectName)
		if err != nil {
			fmt.Println("Remove old new error: ", err)
			return false, err
		}
	}
	var repo string

	if repoURL == "" {
		layout := ""
		prompt := &survey.Select{
			Message: "Please select a layout",
			Options: []string{
				"Advanced",
				"Lite - WIP",
				"Basic - WIP",
				"Chat - WIP",
			},
			Description: func(value string, index int) string {
				if index == 1 {
					return "A lite project structure"
				}
				if index == 2 {
					return "A basic project structure"
				}
				if index == 3 {
					return "A simple chat room containing websocker/tcp"
				}
				return "It has rich functions such as: Wire, Gin, SuaMae, MinhaMae, VossaMae e etc..."
			},
		}
		err := survey.AskOne(prompt, &layout)
		if err != nil {
			return false, err
		}
		if layout == "Advanced" {
			repo = config.RepoFullStructure
		} else if layout == "Chat" {
			repo = config.RepoChat
		} else if layout == "Lite" {
			repo = config.RepoLiteStructure
		} else {
			repo = config.RepoBasicStructure
		}

		err = os.RemoveAll(p.ProjectName)
		if err != nil {
			fmt.Println("remove old new error: ", err)
			return false, err
		}
	} else {
		repo = repoURL
	}

	cmd := exec.Command("git", "clone", repo, p.ProjectName)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("git clone %s error: %s\n", repo, err)
		return false, err
	}
	return true, nil
}

func (p *Project) replaceFiles(packageName string) error {

	err := filepath.Walk(p.ProjectName, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newData := bytes.ReplaceAll(data, []byte(packageName), []byte(p.ProjectName))
		if err := os.WriteFile(path, newData, helper.GetDefaultOSPermissionFile()); err != nil {
			return err
		}
		return nil

	})
	if err != nil {
		fmt.Println("Walk file error: ", err)
		return err
	}
	return nil
}

func (p *Project) replacePackageName() error {
	packageName := helper.GetProjectName(p.ProjectName)

	err := p.replaceFiles(packageName)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "edit", "-module", p.ProjectName)
	cmd.Dir = p.ProjectName
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("go mod edit error: ", err)
		return err
	}
	return nil
}
