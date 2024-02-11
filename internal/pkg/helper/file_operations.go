package helper

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

func getUnixDefaultFilePermissions() os.FileMode {
	return 0644
}

func getWindowsDefaultFilePermissions() os.FileMode {
	return 0666
}

func addLineAfterLastPattern(filename, pattern, newLine string) error {
	unixDefaultFilePermissions := getDefaultOSPermissionFile()

	// Open the file
	file, err := os.OpenFile(filename, os.O_RDWR, unixDefaultFilePermissions)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("TEM ERRO no DEFER")
		}
		log.Printf("Fechou DEFER do file-op do tipo model")
	}(file)

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)
	var lines []string

	// Keep track of whether we found the pattern
	var lastOccurrenceIndex int
	var lastIndentation string

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		// If the current line contains the pattern, update the last occurrence index and capture the indentation
		if strings.Contains(line, pattern) {
			lastOccurrenceIndex = len(lines)
			lastIndentation = extractIndentation(line)
		}
		// Append the current line to the lines slice
		lines = append(lines, line)
	}

	// Insert the new line after the last occurrence index with the same indentation
	lines = insertNewLine(lines, lastOccurrenceIndex+1, lastIndentation, newLine)

	// Write the modified content back to the file
	err = file.Truncate(0)
	if err != nil {
		return err
	} // Clear the file

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	} // Rewind to the beginning

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func getDefaultOSPermissionFile() os.FileMode {
	if runtime.GOOS == "windows" {
		return getWindowsDefaultFilePermissions()
	}
	return getUnixDefaultFilePermissions()
}

func extractIndentation(line string) string {
	var indent string
	for _, char := range line {
		if char == ' ' || char == '\t' {
			indent += string(char)
		} else {
			break
		}
	}
	return indent
}

func insertNewLine(lines []string, index int, indentation, newLine string) []string {
	// If there is no occurrence of the pattern, insert the new line at the end
	if index == 0 {
		return append(lines, newLine)
	}
	// Insert the new line after the last occurrence index with the same indentation
	lines = append(lines[:index], append([]string{indentation + newLine}, lines[index:]...)...)
	return lines
}

func addLineAfterLastPatternWireFile(filename, variableName, newInfo string) error {
	// Open the file
	file, err := os.OpenFile(filename, os.O_RDWR, getDefaultOSPermissionFile())
	log.Printf("Entrou no metodo pra inserir")
	log.Printf(filename, variableName, newInfo)

	if err != nil {
		log.Printf("TEM ERRO")
		return err
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Printf("TEM ERRO no DEFER")

		}
		log.Printf("Fechou DEFER do file-op do tipo wire")
	}(file)

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)
	var lines []string

	// Keep track of whether we found the variable and the last occurrence
	var foundVariable bool
	var foundLastOccurrence bool

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		// If the current line contains the variable, mark it as found
		if strings.Contains(line, variableName) {
			foundVariable = true
		}
		// If we found the variable, look for the last occurrence within it and insert the new info after it
		if foundVariable {
			// If the line contains the pattern and it's the last line inside the variable block, insert the new info
			if strings.Contains(line, ",") && !foundLastOccurrence {
				lines = append(lines, line)
				// Append the new info after the last occurrence within the variable
				lines = append(lines, newInfo)
				foundLastOccurrence = true
			} else {
				// Otherwise, just add the line as is
				lines = append(lines, line)
			}
		} else {
			// Otherwise, just add the line as is
			lines = append(lines, line)
		}
	}

	// Write the modified content back to the file
	err = file.Truncate(0)
	if err != nil {
		return err
	} // Clear the file

	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	} // Rewind to the beginning

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	log.Printf("NAO TEM ERRO")

	return nil
}

func UpdateFile(filetype, filePath, pattern, newLine string) {
	if filetype == "model" {
		_ = addLineAfterLastPattern(filePath, pattern, newLine)
	}
	_ = addLineAfterLastPatternWireFile("source/cmd/server/wire.go", "ServiceSet,", fmt.Sprintf("service.New%sService,", filePath))
}
