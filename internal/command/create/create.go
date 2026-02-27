package create

import (
	"fmt"
	"log"
	"math"
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
	Use:   "create [model|all] [entity-name]",
	Short: "Generate Go source files for an entity",
}

var (
	mustachePath string
	properties   string
)

func init() {
	CmdCreateModel.Flags().StringVarP(&mustachePath, "mustache-path", "t", mustachePath, "Path to custom mustache templates")
	CmdCreateAll.Flags().StringVarP(&mustachePath, "mustache-path", "t", mustachePath, "Path to custom mustache templates")
	CmdCreateModel.Flags().StringVarP(&properties, "properties", "p", "", "Entity fields: title:string;price:float64")
	CmdCreateAll.Flags().StringVarP(&properties, "properties", "p", "", "Entity fields: title:string;price:float64")
}

var CmdCreateModel = &cobra.Command{
	Use:     "model",
	Short:   "Generate a GORM model",
	Example: "gcli create model product title:string price:float64",
	Args:    cobra.MinimumNArgs(1),
	Run:     runCreate,
}

var CmdCreateAll = &cobra.Command{
	Use:     "all",
	Short:   "Generate handler, service, repository, model and router for an entity",
	Example: "gcli create all product title:string price:float64",
	Args:    cobra.MinimumNArgs(1),
	Run:     runCreate,
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

	// Merge --properties flag with Elixir-style positional args (args[1:])
	c.Properties = parseProperties(properties)
	for _, arg := range args[1:] {
		parts := strings.SplitN(arg, ":", 2)
		if len(parts) == 2 {
			name := strutil.UpperFirst(strutil.CamelCase(strings.TrimSpace(parts[0])))
			propType := normalizeGoType(strings.TrimSpace(parts[1]))
			if name != "" && propType != "" {
				c.Properties[name] = propType
			}
		}
	}

	switch c.CreateType {
	case "model":
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

		// Integration test — lives beside the repository file
		origName := c.FileName
		c.CreateType = "test"
		c.FilePath = "source/repository/"
		c.FileName = origName + "_test"
		c.generateFile()
		c.FilePath = ""
		c.FileName = origName
		c.CreateType = "all"

		c.appendToModels()
		c.appendToServerWire()
		c.appendToMigrationWire()
		c.appendToHTTPRouter()

	default:
		log.Fatalf("Invalid create type: %s", c.CreateType)
	}
}

func parseProperties(input string) map[string]string {
	props := make(map[string]string)
	if strings.TrimSpace(input) == "" {
		return props
	}
	// Support both ; and , as separators
	input = strings.ReplaceAll(input, ";", ",")
	for _, pair := range strings.Split(input, ",") {
		parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(parts) == 2 {
			name := strutil.UpperFirst(strutil.CamelCase(strings.TrimSpace(parts[0])))
			propType := normalizeGoType(strings.TrimSpace(parts[1]))
			if name != "" && propType != "" {
				props[name] = propType
			}
		}
	}
	return props
}

// testValueForType returns a hardcoded Go literal usable as a test value for a given Go type.
func testValueForType(goType string) string {
	switch strings.ToLower(goType) {
	case "string":
		return `"test"`
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return "1"
	case "float32", "float64":
		return "1.0"
	case "bool":
		return "true"
	default:
		return "nil"
	}
}

// gormTagForType returns a GORM struct tag for a given Go type.
func gormTagForType(goType string) string {
	switch strings.ToLower(goType) {
	case "string":
		return `gorm:"type:varchar(255);not null"`
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return `gorm:"not null;default:0"`
	case "float32", "float64":
		return `gorm:"type:decimal(10,2);not null;default:0"`
	case "bool":
		return `gorm:"not null;default:false"`
	default:
		return ""
	}
}

// normalizeGoType maps common/Elixir type names to Go types.
func normalizeGoType(t string) string {
	switch strings.ToLower(t) {
	case "integer":
		return "int"
	case "float", "decimal":
		return "float64"
	case "boolean":
		return "bool"
	case "text":
		return "string"
	default:
		return t
	}
}

func (c *Create) generateFile() {
	filePath := c.FilePath
	if filePath == "" {
		filePath = fmt.Sprintf("source/%s/", c.CreateType)
	}

	f := createFile(filePath, strings.ToLower(c.FileName)+".go")
	if f == nil {
		log.Printf("warn: file %s%s already exists, skipping", filePath, strings.ToLower(c.FileName)+".go")
		return
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Printf("warn: failed to close file: %v", err)
		}
	}(f)

	funcMap := template.FuncMap{
		"snake":     strutil.SnakeCase,
		"lower":     strings.ToLower,
		"gormTag":   gormTagForType,
		"testValue": testValueForType,
	}

	var t *template.Template
	var err error
	name := fmt.Sprintf("%s.mustache", c.CreateType)
	if mustachePath == "" {
		t, err = template.New(name).Funcs(funcMap).ParseFS(mustache.CreateTemplateFS, fmt.Sprintf("create/%s", name))
	} else {
		t, err = template.New(name).Funcs(funcMap).ParseFiles(path.Join(mustachePath, name))
	}
	if err != nil {
		log.Fatalf("create %s error: %s", c.CreateType, err.Error())
	}
	if err = t.Execute(f, c); err != nil {
		log.Fatalf("create %s error: %s", c.CreateType, err.Error())
	}

	fileSize, _ := f.Stat()
	kb := math.Round(float64(fileSize.Size()) / 1024)
	log.Printf("Created %s: %s (%.0fkb)", c.CreateType, filePath+strings.ToLower(c.FileName)+".go", kb)
}

func createFile(dirPath string, filename string) *os.File {
	filePath := filepath.Join(dirPath, filename)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Fatalf("Failed to create dir %s: %v", dirPath, err)
	}

	if stat, _ := os.Stat(filePath); stat != nil {
		return nil
	}
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filePath, err)
	}
	return file
}

// appendToModels injects &<StructName>{} into source/model/model.go's
// RetrieveAll() return block so GORM's auto-migrate stays in sync.
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
	indent := "\t\t"

	for i, line := range lines {
		if strings.Contains(line, "return []interface{}{") {
			inBlock = true
			continue
		}
		if inBlock {
			trimmed := strings.TrimSpace(line)
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

// appendToWireSet injects an entry into a named wire.NewSet(...) block.
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

func (c *Create) appendToServerWire() {
	const wireFile = "source/cmd/server/wire.go"
	c.appendToWireSet(wireFile, "repositorySet", fmt.Sprintf("repository.New%sRepository", c.StructName))
	c.appendToWireSet(wireFile, "serviceSet", fmt.Sprintf("service.New%sService", c.StructName))
	c.appendToWireSet(wireFile, "handlerSet", fmt.Sprintf("handler.New%sHandler", c.StructName))
}

func (c *Create) appendToMigrationWire() {
	const wireFile = "source/cmd/migration/wire.go"
	c.appendToWireSet(wireFile, "repositorySet", fmt.Sprintf("repository.New%sRepository", c.StructName))
}

// appendToHTTPRouter injects the handler parameter into NewHTTPServer and
// calls Bind<Entity>Routes(strictAuthRouter, ...) inside the v1 block.
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

	// Inject 1: add handler parameter to NewHTTPServer
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
				if strings.Contains(trimmed, "*handler.") {
					insertAt = i + 1
					indent = line[:len(line)-len(strings.TrimLeft(line, "\t "))]
				}
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

	// Inject 2: Bind<Entity>Routes inside the v1 block
	if !strings.Contains(content, routeCall) {
		lines := strings.Split(content, "\n")
		passedStrictAuth := false
		insertAt := -1

		for i, line := range lines {
			if strings.Contains(line, "strictAuthRouter := v1.Group") {
				passedStrictAuth = true
			}
			if passedStrictAuth && line == "\t}" {
				insertAt = i
				break
			}
		}

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
