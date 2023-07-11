package create

import (
	"fmt"
	"gcli/internal/pkg/helper"
	"gcli/mustache"
	"github.com/spf13/cobra"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Create struct {
	ProjectName        string
	CreateType         string
	FilePath           string
	FileName           string
	FileNameTitleLower string
	FileNameFirstChar  string
	IsFull             bool
}

func NewCreate() *Create {
	return &Create{}
}

var CreateCmd = &cobra.Command{
	Use:     "create [type] [handler-name]",
	Short:   "Create a new handler/service/repository/model",
	Example: "gcli create handler user",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var CreateHandlerCmd = &cobra.Command{
	Use:     "handler",
	Short:   "Create a new ",
	Example: "gcli create handler user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}

var CreateServiceCmd = &cobra.Command{
	Use:     "service",
	Short:   "Create a new ",
	Example: "gcli create service user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}

var CreateRepositoryCmd = &cobra.Command{
	Use:     "repository",
	Short:   "Create a new ",
	Example: "gcli create repository user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}

var CreateModelCmd = &cobra.Command{
	Use:     "model",
	Short:   "Create a new ",
	Example: "gcli create model user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}

var CreateAllCmd = &cobra.Command{
	Use:     "all",
	Short:   "Create a new ",
	Example: "gcli create all user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	c := NewCreate()
	c.ProjectName = helper.GetProjectName(".")
	c.CreateType = cmd.Use
	c.FilePath, c.FileName = filepath.Split(args[0])
	c.FileName = strings.ReplaceAll(strings.ToUpper(string(c.FileName[0]))+c.FileName[1:], ".go", "")
	c.FileNameTitleLower = strings.ToLower(string(c.FileName[0])) + c.FileName[1:]
	c.FileNameFirstChar = string(c.FileNameTitleLower[0])

	switch c.CreateType {
	case "handler", "service", "repository", "model":
		c.generateFile()
	case "all":
		c.CreateType = "handler"
		c.generateFile()

		c.CreateType = "service"
		c.generateFile()

		c.CreateType = "repository"
		c.generateFile()

		c.CreateType = "model"
		c.generateFile()
	default:
		log.Fatalf("Invalid handler type %s", c.CreateType)
	}
}

func (c *Create) generateFile() {
	filePath := c.FilePath
	if filePath == "" {
		filePath = fmt.Sprintf("source/%s/", c.CreateType)
	}

	f := createFile(filePath, strings.ToLower(c.FileName)+".go")
	if f == nil {
		log.Printf("warn: file %s%s %s", filePath, strings.ToLower(c.FileName)+".go", "already exists")
		return
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	t, err := template.ParseFS(mustache.CreateTemplateFS, fmt.Sprintf("create/%s.mustache", c.CreateType))
	if err != nil {
		log.Fatalf("create %s error: %s", c.CreateType, err.Error())
	}
	err = t.Execute(f, c)
	if err != nil {
		log.Fatalf("create %s error: %s", c.CreateType, err.Error())
	}
	fileSize, _ := f.Stat()

	kilobytes := math.Round(float64(fileSize.Size()) / 1024)

	log.Printf("Created new %s: %s (%vkb)", c.CreateType, filePath+strings.ToLower(c.FileName)+".go", kilobytes)
}

func createFile(dirPath string, filename string) *os.File {
	filePath := dirPath + filename

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create dir %s: %v", dirPath, err)
	}
	stat, _ := os.Stat(filePath)
	if stat != nil {
		return nil
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filePath, err)
	}

	return file
}
