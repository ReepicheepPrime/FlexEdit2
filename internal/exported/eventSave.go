package exported

import (
	"strings"
	
	"github.com/reepicheepprime/flexible-survival-editor/helpers"
)

const EventSaveFile = "FSEventSave.glkdata"

type GameEvents []*GameEvent

type GameEvent struct {
	Name          string
	ResolveState  string // ["Resolved", "Unresolved"]
	ActiveState   string // ["Active", "Inactive"]
	Resolution    int
	SituationArea string
}

func ParseEventSave(raw string) (GameEvents, error) {
	events := GameEvents{}
	rawLines := strings.Split(raw, "\n")[2:] // skip header

	for _, rawLine := range rawLines {
		rawValues := strings.Split(rawLine, " ")

		// parse Name
		Name, err := helpers.DecodeString(rawValues[0])
		if err != nil {
			return nil, err
		}

		// parse ResolveState
		ResolveState, err := helpers.DecodeString(rawValues[1])
		if err != nil {
			return nil, err
		}

		// parse ActiveState
		ActiveState, err := helpers.DecodeString(rawValues[2])
		if err != nil {
			return nil, err
		}

		// parse Resolution
		Resolution, err := helpers.DecodeNumber(rawValues[3])
		if err != nil {
			return nil, err
		}

		// parse SituationArea
		SituationArea, err := helpers.DecodeString(rawValues[2])
		if err != nil {
			return nil, err
		}

		event := &GameEvent{Name, ResolveState, ActiveState, Resolution, SituationArea}
		events = append(events, event)
	}

	return events, nil
}
