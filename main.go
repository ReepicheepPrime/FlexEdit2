package main

import (
	"os"
	"path"
	"strings"

	"github.com/reepicheepprime/flexible-survival-editor/internal/exported"
)

func loadRawFile(path string) string {
	rawFile, err := os.ReadFile(path)
	if err != nil { panic(err) }

	decoded := string(rawFile)
	return strings.ReplaceAll(decoded, "\r\n", "\n") // handle windows-generated files
}

func parseExported(baseDir string) {
	// FSEventSave
	rawEventSave := loadRawFile(path.Join(baseDir, exported.EventSaveFile))
	eventSave, _ := exported.ParseEventSave(rawEventSave)
	panic(eventSave)
}

func main() {
	baseDir := "testdata"
	parseExported(baseDir)
}