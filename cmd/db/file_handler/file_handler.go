package file_handler

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	dbutils "github.com/AndreDrummer/school-system-api/cmd/db/utils"
	apperrors "github.com/AndreDrummer/school-system-api/cmd/errors"

	"strings"
)

func OpenFileWithPerm(filename string, flag int) (*os.File, error) {
	file, err := os.OpenFile(filename, flag, 0644)

	if err != nil {
		log.Fatalf("ERROR %v opening file %v", err, filename)
		return nil, err
	}

	return file, nil
}

func PrintFileContent(file *os.File) {
	fileContent := GetFileContent(file)

	for _, v := range fileContent {
		fmt.Println(v)
	}
}

func GetFileContent(file *os.File) []string {
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	var content []string

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			content = append(content, line)
		}
	}
	return content
}

func GetFileEntryByPrefix(file *os.File, prefix string) (string, error) {
	fileContent := GetFileContent(file)

	for _, v := range fileContent {
		vPrefix := strings.Split(v, " ")[0]
		if vPrefix == prefix {
			return v, nil
		}
	}

	return "", errors.New("entry not found")
}

func UpdateFileEntry(file *os.File, entryPrefix, updatedEntry string) error {
	fileContent := GetFileContent(file)
	var newContent []string
	var found bool

	for _, v := range fileContent {

		vPrefix := strings.Split(v, " ")[0]
		if vPrefix == entryPrefix {
			found = true
			newContent = append(newContent, updatedEntry)
		} else if v == "" {
			continue
		} else {
			newContent = append(newContent, v)
		}
	}

	if found {
		OverrideFileContent(file, newContent)
		return nil
	} else {
		return &apperrors.NotFoundError{}
	}
}

func RemoveFileEntry(file *os.File, entryPrefix any) {
	fileContent := GetFileContent(file)
	var newContent []string

	for _, v := range fileContent {
		vPrefix := strings.Split(v, " ")[0]
		if vPrefix == entryPrefix || v == "" {
			continue
		} else {
			newContent = append(newContent, v)
		}
	}

	OverrideFileContent(file, newContent)
}

func OverrideFileContent(file *os.File, content []string) {
	file.Truncate(0)

	file.Seek(0, 0)

	dbutils.SortSliceStringByID(content, " ")

	for _, v := range content {
		file.WriteString(fmt.Sprintf("%s\n", v))
	}
}

func AppendToFile(file *os.File, content string) {
	file.WriteString(fmt.Sprintf("\n%s", content))
}

func ClearFileContent(file *os.File) {
	file.Truncate(0)
}
