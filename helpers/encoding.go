package helpers

import (
	"fmt"
	"strconv"
	"strings"
)

func DecodeString(s string) (string, error) {
	if !strings.HasPrefix(s, "S") {
		return "", fmt.Errorf("input value is not a string: %s", s)
	}

	s = strings.TrimPrefix(s, "S")
	s = strings.TrimSuffix(s, ";")

	out := ""
	chars := strings.Split(s, ",")

	for _, char := range chars {
		num, err := strconv.Atoi(char)
		if err != nil {
			return "", fmt.Errorf("failed to convert to int: %s", char)
		}
		out += string(rune(num))
	}

	return out, nil
}

func DecodeNumber(s string) (int, error) {
	return strconv.Atoi(strings.TrimSuffix(s, ";"))
}