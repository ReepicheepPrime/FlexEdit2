package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type GameVariableID struct {
	Name    string
	Type    string
	RawType string
	RawName string
}

func parseType(rawType string) (string, error) {
	out := ""

	// do NOT handle recursive slices
	if strings.Contains(rawType, "list of lists of") {
		return "", fmt.Errorf("recursive slice parsing is not supported")
	}

	// handle slices
	if strings.HasPrefix(rawType, "list of") {
		rawType = strings.ReplaceAll(rawType, "list of ", "")
		rawType = strings.TrimSuffix(rawType, "s")
		out += "[]"
	}

	// convert to Go primatives
	if rawType == "number" {
		return out + "int", nil
	}

	if rawType == "text" {
		return out + "string", nil
	}

	if rawType == "truth state" {
		return out + "bool", nil
	}

	return "", fmt.Errorf("unsupported type: %s", rawType)
}

func getGameTables() string {
	baseUrl := "https://raw.githubusercontent.com/Nuku/Flexible-Survival"
	version := "master"
	fileUrl := "Core%20Mechanics/GameTables.i7x"

	// init request
	res, err := http.Get(fmt.Sprintf("%s/%s/%s", baseUrl, version, fileUrl))
	if err != nil {
		log.Fatal("Failed to request GameTables.i7x from github")
	}

	// read body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Failed to read GameTables.i7x from HTTP response")
	}

	return string(body)
}

func parseGameVariableIDs(rawGameVariableIDs string) []GameVariableID {
	var gameVars []GameVariableID
	gameVarPairs := strings.Split(rawGameVariableIDs, "\n")

	for _, gameVarPair := range gameVarPairs {
		pattern := regexp.MustCompile(`\"([\w -]+)\"\t\"([\w ]+)\"`)
		gameVar := pattern.FindStringSubmatch(gameVarPair)[1:]

		goType, err := parseType(gameVar[1])
		if err != nil {
			log.Printf("failed to parse type: %s", err.Error())
			continue
		}

		cleanName := regexp.MustCompile(`\W`).ReplaceAllString(gameVar[0], "_")

		gameVars = append(gameVars, GameVariableID{
			Name:    strings.Title(cleanName),
			Type:    goType,
			RawName: gameVar[0],
			RawType: gameVar[1],
		})
	}

	return gameVars
}

// RightPad2Len https://github.com/DaddyOh/golang-samples/blob/master/pad.go
func RightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

func main() {
	rawGameTables := getGameTables()

	// extract only the block we are looking for
	headerGameVarsWithExcess := strings.Split(rawGameTables, "Table of GameVariableIDs\n")[1]
	headerGameVars := strings.Split(headerGameVarsWithExcess, "\n\n")[0]
	rawGameVars := strings.ReplaceAll(headerGameVars, "Name(text)	Type(text)\n", "")

	// parse the individual game vars
	gameVars := parseGameVariableIDs(rawGameVars)

	// find longest name
	longestName := 0
	longestRawName := 0
	for _, gameVar := range gameVars {
		if len(gameVar.Name) > longestName {
			longestName = len(gameVar.Name)
		}
		if len(gameVar.RawName) > longestRawName {
			longestRawName = len(gameVar.RawName)
		}
	}

	// longest type is []number and there for hardcoded for now
	longestType := 8

	// generate go file
	output := "// File generated automatically. DO NOT EDIT.\npackage main"
	output += "\n\ntype GameVariables struct {"
	for _, gameVar := range gameVars {
		// add custom tags to help maintain consitency with FS. This mostly only exists because some dick
		// used a hyphen in a variable name, which Go doesn't support for struct value names.
		jsonTag := RightPad2Len(fmt.Sprintf("json:\"%s\"", gameVar.RawName), " ", longestRawName)
		i7NameTag := RightPad2Len(fmt.Sprintf("i7Name:\"%s\"", gameVar.RawName), " ", longestRawName)
		i7TypeTag := fmt.Sprintf("i7Type:\"%s\"", gameVar.RawType)
		fullTags := strings.Join([]string{jsonTag, i7NameTag, i7TypeTag}, " ")

		paddedName := RightPad2Len(gameVar.Name, " ", longestName)
		paddedType := RightPad2Len(gameVar.Type, " ", longestType)
		output += fmt.Sprintf("\n\t %s %s `%s`", paddedName, paddedType, fullTags)
	}
	output += "\n}"

	// write output to file
	outFile, _ := filepath.Abs("internal/gameTables.go")
	err := os.WriteFile(outFile, []byte(output), 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to write file: [%s]", err))
	}
}
