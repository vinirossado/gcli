package create

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/spf13/cobra"

	"github.com/vinirossado/gcli/internal/pkg/helper"
	"github.com/vinirossado/gcli/mustache"
)

type Create struct {
	ProjectName          string
	CreateType           string
	FilePath             string
	FileName             string
	StructName           string
	StructNameLowerFirst string
	StructNameFirstChar  string
	StructNameSnakeCase  string
	IsFull               bool
	Properties           map[string]string
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
	Short:   "Create a new handler ",
	Example: "gcli create handler user",
	Args:    cobra.MinimumNArgs(1),
	Run:     runCreate,
}

var CmdCreateService = &cobra.Command{
	Use:     "service",
	Short:   "Create a new service",
	Example: "gcli create service user",
	Args:    cobra.MinimumNArgs(1),
	Run:     runCreate,
}

var CmdCreateRepository = &cobra.Command{
	Use:     "repository",
	Short:   "Create a new repository ",
	Example: "gcli create repository user",
	Args:    cobra.MinimumNArgs(1),
	Run:     runCreate,
}

var CmdCreateModel = &cobra.Command{
	Use:     "model",
	Short:   "Create a new model ",
	Example: "gcli create model user",
	Args:    cobra.MinimumNArgs(1),
	Run:     runCreate,
}

var CmdCreateAll = &cobra.Command{
	Use:     "all",
	Short:   "Create a new handler & service & repository & model ",
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

	/*
		c := NewCreate()
		c.ProjectName = helper.GetProjectName(".")
		c.CreateType = cmd.Use
		c.FilePath, c.StructName = filepath.Split(args[0])
		c.FileName = strings.ReplaceAll(c.StructName, ".go", "")
		c.StructName = strutil.UpperFirst(strutil.CamelCase(c.FileName))
		c.StructNameLowerFirst = strutil.LowerFirst(c.StructName)
		c.StructNameFirstChar = string(c.StructNameLowerFirst[0])
		c.StructNameSnakeCase = strutil.SnakeCase(c.StructName)

	*/
	c := NewCreate()
	c.ProjectName = helper.GetProjectName(".")
	c.CreateType = cmd.Use
	c.FilePath, c.StructName = filepath.Split(args[0])
	c.FileName = strings.ReplaceAll(c.StructName, ".go", "")
	c.StructName = strutil.UpperFirst(strutil.CamelCase(c.FileName))
	c.StructNameLowerFirst = strutil.LowerFirst(c.StructName)
	c.StructNameFirstChar = string(c.StructNameLowerFirst[0])
	c.StructNameSnakeCase = strutil.SnakeCase(c.StructName)
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
	// if strings.Contains(helper.GetProjectRootName(), "gcli") {
	// 	filePath = fmt.Sprintf("Debug/source/%s/", c.CreateType)
	// }

	if filePath == "" {
		filePath = fmt.Sprintf("source/%s/", c.CreateType)
	}

	f := createFile(filePath, strings.ToLower(c.FileName)+".go")
	if f == nil {
		log.Printf("warn: file %s%s %s", filePath, strings.ToLower(c.FileName)+".go", "already exists")
		return
	}
	defer func(f *os.File) {
		// 	helper.UpdateFile(c.CreateType, filePath, "{},", fmt.Sprintf("&%s{},", c.FileName))
		f.Close()
	}(f)

	var t *template.Template
	var err error
	if mustachePath == "" {
		t, err = template.ParseFS(mustache.CreateTemplateFS, fmt.Sprintf("create/%s.mustache", c.CreateType))
	} else {
		t, err = template.ParseFiles(path.Join(mustachePath, fmt.Sprintf("%s.mustache", c.CreateType)))
	}
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

	filePath := filepath.Join(dirPath, filename)

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

// func containsFilterString(line string) bool {
// 	print(line)
// 	if strings.Contains(line, "return r") {
// 		return true
// 	}
// 	return false
// }
