package create

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vinirossado/gcli/internal/pkg/helper"
	"github.com/vinirossado/gcli/mustache"
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

	//updateFile(filePath, strings.ToLower(c.FileName)+".go")
}

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

func updateFile(dirPath string, filename string) {
	// Open the file for reading
	fileName := filename
	print(dirPath + "http.go")

	file, err := os.Open(dirPath + fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Create a temporary buffer to store the filtered lines and the new specific line
	tempBuffer := ""

	// Filter lines and append the new specific line if necessary
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Filter lines based on your criteria, for example, lines containing "filter_string"
		if containsFilterString(line) {
			tempBuffer += line + "\n" // Add the selected line
		} else {
			tempBuffer += line + "\n"            // Add other lines as they are
			tempBuffer += "Addicionou esse k7\n" // Add the new specific line
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return
	}

	// Reopen the file for writing
	outputFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer func(outputFile *os.File) {
		err := outputFile.Close()
		if err != nil {

		}
	}(outputFile)

	// Write the contents of the temporary buffer (including selected lines and new specific line) to the file
	_, err = outputFile.WriteString(tempBuffer)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("File filtering and writing completed successfully.")
}

func containsFilterString(line string) bool {
	print(line)
	if strings.Contains(line, "return r") {
		return true
	}
	return false
}
