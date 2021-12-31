package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/reepicheepprime/flexible-survival-editor/internal/exported"
)

func loadRawFile(path string) string {
	rawFile, err := os.ReadFile(path)
	if err != nil { panic(err) }

	decoded := string(rawFile)
	unixed := strings.ReplaceAll(decoded, "\r\n", "\n") // handle windows-generated files
	return strings.TrimSuffix(unixed, "\n") // remove ending newline
}

func parseExported(baseDir string) {
	// FSEventSave
	rawEventSave := loadRawFile(path.Join(baseDir, exported.EventSaveFile))
	eventSave, err := exported.ParseEventSave(rawEventSave)
	if err != nil {
		panic(err)
	}
	for _, event := range eventSave {
		fmt.Printf("%+v\n", *event)
	}
}

func main() {
	baseDir := "testdata"
	parseExported(baseDir)
}