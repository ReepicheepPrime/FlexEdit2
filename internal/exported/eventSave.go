package exported

import (
	"fmt"
	"strconv"
	"strings"
)

const EventSaveFile = "FSEventSave.glkdata"

type GameEvents []*GameEvent

type GameEvent struct {
	Name          string
	ResolveState  string
	ActiveState   string
	Resolution    int
	SituationArea string
}

func DecodeAscii(s string) (string, error) {
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

func ParseEventSave(raw string) (GameEvents, error) {
	events := GameEvents{}
	rawLines := strings.Split(raw, "\n")[2:] // skip header

	for _, rawLine := range rawLines {
		rawValues := strings.Split(rawLine, " ")

		// parse Name
		Name, err := DecodeAscii(rawValues[0])
		if err != nil {
			return nil, err
		}

		// parse ResolveState
		ResolveState, err := DecodeAscii(rawValues[1])
		if err != nil {
			return nil, err
		}

		// parse ActiveState
		ActiveState, err := DecodeAscii(rawValues[2])
		if err != nil {
			return nil, err
		}

		// parse Resolution
		Resolution, err := strconv.Atoi(strings.TrimSuffix(rawValues[3], ";"))
		if err != nil {
			return nil, err
		}

		// parse SituationArea
		SituationArea, err := DecodeAscii(rawValues[2])
		if err != nil {
			return nil, err
		}

		event := &GameEvent{Name, ResolveState, ActiveState, Resolution, SituationArea}
		events = append(events, event)
	}

	return events, nil
}
