package events

import (
	"encoding"
	"fmt"
)

// RecordingStatus is an Enum indicating the state of the conference call recording.
type RecordingStatus string

var _ Validator = (*RecordingStatus)(nil)
var _ encoding.TextUnmarshaler = (*RecordingStatus)(nil)

// RecordingStatus Enumerations
const (
	RecordingInProgress RecordingStatus = "in-progress"
	RecordingCompleted  RecordingStatus = "completed"
	RecordingAbsent     RecordingStatus = "absent"
)

// Validate checks that the current value matches one of the allowed enums.
func (es RecordingStatus) Validate() error {
	switch es {
	case RecordingInProgress,
		RecordingCompleted,
		RecordingAbsent:
		return nil
	}
	return fmt.Errorf("unknown %T value: %s", es, es)
}

// UnmarshalText implements the TextUnmarshaler interface.
func (es *RecordingStatus) UnmarshalText(text []byte) (err error) {
	check := RecordingStatus(text)

	if err := check.Validate(); err != nil {
		return err
	}

	*es = check

	return nil
}

type Recording struct {
	AccountSID         string `json:"AccountSid"`
	ConferenceSID      string `json:"ConferenceSid"`
	RecordingSID       string `json:"RecordingSid"`
	RecordingURL       string `json:"RecordingUrl"`
	RecordingStatus    RecordingStatus
	RecordingDuration  int
	RecordingChannels  int
	RecordingStartTime TimeRFC1123Z
	RecordingSource    string
}

var _ Validator = (*Recording)(nil)

// Validate checks that the current value matches one of the allowed enums.
func (c Recording) Validate() error {
	switch {
	case !IsValid(c.RecordingStatus):
		return fmt.Errorf("event RecordingStatus invalid got='%s'", c.RecordingStatus)
	}

	return nil
}
