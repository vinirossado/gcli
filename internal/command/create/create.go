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

	"github.com/spf13/cobra"
	"github.com/vinirossado/strutil"

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

		c.appendToModels()
		c.appendToServerWire()
		c.appendToMigrationWire()
		c.appendToHTTPRouter()

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

// appendToModels finds source/model/model.go and injects &<StructName>{} into
// the return []interface{}{} block, so GORM's auto-migrate stays in sync.
func (c *Create) appendToModels() {
	const modelsFile = "source/model/model.go"

	data, err := os.ReadFile(modelsFile)
	if err != nil {
		log.Printf("warn: %s not found, skipping model registration", modelsFile)
		return
	}

	entry := fmt.Sprintf("&%s{}", c.StructName)
	if strings.Contains(string(data), entry) {
		log.Printf("warn: %s is already registered in %s", entry, modelsFile)
		return
	}

	lines := strings.Split(string(data), "\n")
	inBlock := false
	insertAt := -1
	indent := "\t\t" // default; overwritten when we find an existing entry

	for i, line := range lines {
		if strings.Contains(line, "return []interface{}{") {
			inBlock = true
			continue
		}
		if inBlock {
			trimmed := strings.TrimSpace(line)
			// Capture indentation from the first existing entry (e.g. "\t\t&User{},")
			if strings.HasPrefix(trimmed, "&") {
				indent = line[:len(line)-len(strings.TrimLeft(line, "\t "))]
			}
			if trimmed == "}" {
				insertAt = i
				break
			}
		}
	}

	if insertAt == -1 {
		log.Printf("warn: could not find 'return []interface{}{}' block in %s", modelsFile)
		return
	}

	newLine := fmt.Sprintf("%s%s,", indent, entry)
	result := make([]string, 0, len(lines)+1)
	result = append(result, lines[:insertAt]...)
	result = append(result, newLine)
	result = append(result, lines[insertAt:]...)

	if err := os.WriteFile(modelsFile, []byte(strings.Join(result, "\n")), helper.GetDefaultOSPermissionFile()); err != nil {
		log.Printf("warn: failed to update %s: %v", modelsFile, err)
		return
	}

	log.Printf("Registered %s in %s", entry, modelsFile)
}

// appendToWireSet injects an entry into a named wire.NewSet(...) block inside filePath.
func (c *Create) appendToWireSet(filePath, setVarName, entry string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("warn: %s not found, skipping", filePath)
		return
	}

	if strings.Contains(string(data), entry) {
		log.Printf("warn: %s already in %s", entry, filePath)
		return
	}

	anchor := fmt.Sprintf("var %s = wire.NewSet(", setVarName)
	lines := strings.Split(string(data), "\n")
	inBlock := false
	insertAt := -1
	indent := "\t"

	for i, line := range lines {
		if strings.Contains(line, anchor) {
			inBlock = true
			continue
		}
		if inBlock {
			trimmed := strings.TrimSpace(line)
			if trimmed == ")" {
				insertAt = i
				break
			}
			if len(trimmed) > 0 {
				indent = line[:len(line)-len(strings.TrimLeft(line, "\t "))]
			}
		}
	}

	if insertAt == -1 {
		log.Printf("warn: could not find '%s' block in %s", anchor, filePath)
		return
	}

	newLine := fmt.Sprintf("%s%s,", indent, entry)
	result := make([]string, 0, len(lines)+1)
	result = append(result, lines[:insertAt]...)
	result = append(result, newLine)
	result = append(result, lines[insertAt:]...)

	if err := os.WriteFile(filePath, []byte(strings.Join(result, "\n")), helper.GetDefaultOSPermissionFile()); err != nil {
		log.Printf("warn: failed to update %s: %v", filePath, err)
		return
	}
	log.Printf("Registered %s in %s (%s)", entry, filePath, setVarName)
}

// appendToServerWire injects the repository, service and handler constructors
// into source/cmd/server/wire.go.
func (c *Create) appendToServerWire() {
	const wireFile = "source/cmd/server/wire.go"
	c.appendToWireSet(wireFile, "repositorySet", fmt.Sprintf("repository.New%sRepository", c.StructName))
	c.appendToWireSet(wireFile, "serviceSet", fmt.Sprintf("service.New%sService", c.StructName))
	c.appendToWireSet(wireFile, "handlerSet", fmt.Sprintf("handler.New%sHandler", c.StructName))
}

// appendToMigrationWire injects the repository constructor into
// source/cmd/migration/wire.go.
func (c *Create) appendToMigrationWire() {
	const wireFile = "source/cmd/migration/wire.go"
	c.appendToWireSet(wireFile, "repositorySet", fmt.Sprintf("repository.New%sRepository", c.StructName))
}

// appendToHTTPRouter injects the new handler into NewHTTPServer's parameter list
// and calls Bind<Entity>Routes(strictAuthRouter, ...) inside the v1 block so
// routes land under /v1/<entity> with auth middleware already applied.
func (c *Create) appendToHTTPRouter() {
	const httpFile = "source/router/http.go"

	data, err := os.ReadFile(httpFile)
	if err != nil {
		log.Printf("warn: %s not found, skipping", httpFile)
		return
	}

	handlerParam := fmt.Sprintf("%sHandler *handler.%sHandler,", c.StructNameLowerFirst, c.StructName)
	routeCall := fmt.Sprintf("Bind%sRoutes(strictAuthRouter, jwt, *%sHandler, logger)", c.StructName, c.StructNameLowerFirst)
	content := string(data)

	// --- Inject 1: add handler parameter to NewHTTPServer ---
	if !strings.Contains(content, handlerParam) {
		lines := strings.Split(content, "\n")
		inFunc := false
		insertAt := -1
		indent := "\t"

		for i, line := range lines {
			if strings.Contains(line, "func NewHTTPServer(") {
				inFunc = true
				continue
			}
			if inFunc {
				trimmed := strings.TrimSpace(line)
				// Track the last *handler.XxxHandler parameter line
				if strings.Contains(trimmed, "*handler.") {
					insertAt = i + 1
					indent = line[:len(line)-len(strings.TrimLeft(line, "\t "))]
				}
				// Stop at the closing ) of the parameter list
				if strings.HasPrefix(trimmed, ")") {
					break
				}
			}
		}

		if insertAt != -1 {
			newLine := fmt.Sprintf("%s%s", indent, handlerParam)
			result := make([]string, 0, len(lines)+1)
			result = append(result, lines[:insertAt]...)
			result = append(result, newLine)
			result = append(result, lines[insertAt:]...)
			content = strings.Join(result, "\n")
			log.Printf("Added %s parameter to NewHTTPServer in %s", handlerParam, httpFile)
		} else {
			log.Printf("warn: could not find handler parameter location in %s", httpFile)
		}
	}

	// --- Inject 2: Bind<Entity>Routes(strictAuthRouter, ...) inside the v1 block ---
	// We look for the single-tab "}" that closes the v1 outer block and insert
	// just before it, where strictAuthRouter is still in scope.
	if !strings.Contains(content, routeCall) {
		lines := strings.Split(content, "\n")
		passedStrictAuth := false
		insertAt := -1

		for i, line := range lines {
			if strings.Contains(line, "strictAuthRouter := v1.Group") {
				passedStrictAuth = true
			}
			// The v1 outer block closes at a single-tab "}" — exactly "\t}"
			if passedStrictAuth && line == "\t}" {
				insertAt = i
				break
			}
		}

		// Fallback: inject before "return s" if the v1 pattern wasn't found
		if insertAt == -1 {
			for i, line := range lines {
				if strings.TrimSpace(line) == "return s" {
					insertAt = i
					break
				}
			}
		}

		if insertAt != -1 {
			newLine := fmt.Sprintf("\t\t%s", routeCall)
			result := make([]string, 0, len(lines)+1)
			result = append(result, lines[:insertAt]...)
			result = append(result, newLine)
			result = append(result, lines[insertAt:]...)
			content = strings.Join(result, "\n")
			log.Printf("Added %s inside v1 block in %s", routeCall, httpFile)
		} else {
			log.Printf("warn: could not find injection point in %s", httpFile)
		}
	}

	if err := os.WriteFile(httpFile, []byte(content), helper.GetDefaultOSPermissionFile()); err != nil {
		log.Printf("warn: failed to update %s: %v", httpFile, err)
	}
}

// func containsFilterString(line string) bool {
// 	print(line)
// 	if strings.Contains(line, "return r") {
// 		return true
// 	}
// 	return false
// }
