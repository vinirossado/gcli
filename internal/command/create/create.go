package create

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vinirossado/gcli/internal/pkg/helper"
	"github.com/vinirossado/gcli/mustache"
	"log"
	"math"
	"math/rand"
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
	Properties         map[string]string
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

var (
	mustachePath string
)

func init() {
	CmdCreateHandler.Flags().StringVarP(&mustachePath, "mustache-path", "t", mustachePath, "mustache path")
	CmdCreateService.Flags().StringVarP(&mustachePath, "mustache-path", "t", mustachePath, "mustache path")
	CmdCreateRepository.Flags().StringVarP(&mustachePath, "mustache-path", "t", mustachePath, "mustache path")
	CmdCreateModel.Flags().StringVarP(&mustachePath, "mustache-path", "t", mustachePath, "mustache path")
	CmdCreateAll.Flags().StringVarP(&mustachePath, "mustache-path", "t", mustachePath, "mustache path")
	CmdCreateModel.Flags().StringVarP(&properties, "properties", "p", "", "Properties of the model entity (format: name:type)")
}

var CmdCreateHandler = &cobra.Command{
	Use:     "handler",
	Short:   "Create a new ",
	Example: "gcli create handler user",
	Args:    cobra.MinimumNArgs(1),
	Run:     runCreate,
}

var CmdCreateService = &cobra.Command{
	Use:     "service",
	Short:   "Create a new ",
	Example: "gcli create service user",
	Args:    cobra.MinimumNArgs(1),
	Run:     runCreate,
}

var CmdCreateRepository = &cobra.Command{
	Use:     "repository",
	Short:   "Create a new ",
	Example: "gcli create repository user",
	Args:    cobra.MinimumNArgs(1),
	Run:     runCreate,
}

var CmdCreateModel = &cobra.Command{
	Use:     "model",
	Short:   "Create a new ",
	Example: "gcli create model user",
	Args:    cobra.MinimumNArgs(1),
	Run:     runCreate,
}

var CmdCreateAll = &cobra.Command{
	Use:     "all",
	Short:   "Create a new ",
	Example: "gcli create all user",
	Args:    cobra.MinimumNArgs(1),
	Run:     runCreate,
}
var properties string

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func runCreate(cmd *cobra.Command, args []string) {
	c := NewCreate()
	c.ProjectName = helper.GetProjectName(".")
	c.CreateType = cmd.Use
	c.FilePath, c.FileName = filepath.Split(args[0])
	//c.FileName = strings.ReplaceAll(strings.ToUpper(string(c.FileName[0]))+c.FileName[1:], ".go", "")
	c.FileName = RandStringBytes(10) + ".go"
	c.FileNameTitleLower = strings.ToLower(string(c.FileName[0])) + c.FileName[1:]
	c.FileNameFirstChar = string(c.FileNameTitleLower[0])
	c.Properties = parseProperties(properties)

	switch c.CreateType {

	case "handler", "service", "repository", "model", "router":
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

		c.CreateType = "router"
		c.generateFile()

	default:
		log.Fatalf("Invalid handler type %s", c.CreateType)
	}
}

func parseProperties(properties string) map[string]string {
	props := make(map[string]string)

	pairs := strings.Split(properties, ",")

	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) == 2 {
			propName := strings.TrimSpace(parts[0])
			propType := strings.TrimSpace(parts[1])
			props[propName] = propType
		}
	}

	return props
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
		log.Printf(c.CreateType)
		helper.UpdateFile(c.CreateType, filePath, "{},", fmt.Sprintf("&%s{},", c.FileName))
		f.Close()
		log.Printf("Fechou DEFER do Generate")
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

// TODO: Rename Method
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

func containsFilterString(line string) bool {
	print(line)
	if strings.Contains(line, "return r") {
		return true
	}
	return false
}
