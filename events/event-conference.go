package events

import (
	"encoding"
	"fmt"
)

/*

All Events:
	ConferenceSid
	FriendlyName
	AccountSid
	SequenceNumber
	Timestamp
	StatusCallbackEvent

join leave mute hold speaker: (type ConfrenceCall)
	CallSid
	Muted
	Hold
	Coaching
	EndConferenceOnExit
	StartConferenceOnEnter

end: (type ConfrenceCallEnd)
	CallSidEndingConference
	ParticipantLabelEndingConference
	ReasonConferenceEnded
	Reason

announcement: (type ConfrenceCallAnnouncement)
	CallSid
	Muted
	Hold
	Coaching
	EndConferenceOnExit
	StartConferenceOnEnter
	ReasonAnnouncementFailed
	AnnounceUrl

*/

// ConferenceStatus is an Enum indicating the state of the conference call.
type ConferenceStatus string

var _ Validator = (*ConferenceStatus)(nil)
var _ encoding.TextUnmarshaler = (*ConferenceStatus)(nil)

// ConferenceStatus Enumerations
const (
	StatusConferenceEnd          ConferenceStatus = "conference-end"
	StatusConferenceStart        ConferenceStatus = "conference-start"
	StatusParticipantLeave       ConferenceStatus = "participant-leave"
	StatusParticipantJoin        ConferenceStatus = "participant-join"
	StatusParticipantMute        ConferenceStatus = "participant-mute"
	StatusParticipantUnmute      ConferenceStatus = "participant-unmute"
	StatusParticipantHold        ConferenceStatus = "participant-hold"
	StatusParticipantUnhold      ConferenceStatus = "participant-unhold"
	StatusParticipantSpeechStart ConferenceStatus = "participant-speech-start"
	StatusParticipantSpeechStop  ConferenceStatus = "participant-speech-stop"
	StatusAnnouncementEnd        ConferenceStatus = "announcement-end"
	StatusAnnouncementFail       ConferenceStatus = "announcement-fail"
)

// Validate checks that the current value matches one of the allowed enums.
func (es ConferenceStatus) Validate() error {
	switch es {
	case StatusConferenceEnd,
		StatusConferenceStart,
		StatusParticipantLeave,
		StatusParticipantJoin,
		StatusParticipantMute,
		StatusParticipantUnmute,
		StatusParticipantHold,
		StatusParticipantUnhold,
		StatusParticipantSpeechStart,
		StatusParticipantSpeechStop,
		StatusAnnouncementEnd,
		StatusAnnouncementFail:
		return nil
	}
	return fmt.Errorf("unknown %T value: %s", es, es)
}

// UnmarshalText implements the TextUnmarshaler interface.
func (es *ConferenceStatus) UnmarshalText(text []byte) (err error) {
	check := ConferenceStatus(text)

	if err := check.Validate(); err != nil {
		return err
	}

	*es = check

	return nil
}

// ReasonConferenceEnded is an Enum indicating the reason the conference call ended.
type ReasonConferenceEnded string

var _ Validator = (*ReasonConferenceEnded)(nil)
var _ encoding.TextUnmarshaler = (*ConferenceStatus)(nil)

const (
	ReasonEndedViaAPI                    ReasonConferenceEnded = "conference-ended-via-api"
	ReasonParticipantKicked              ReasonConferenceEnded = "last-participant-kicked"
	ReasonParticipantLeft                ReasonConferenceEnded = "last-participant-left"
	ReasonParticipantKickedWithEndOnExit ReasonConferenceEnded = "participant-with-end-conference-on-exit-kicked"
	ReasonParticipantLeftWithEndOnExit   ReasonConferenceEnded = "participant-with-end-conference-on-exit-left"
)

// Validate checks that the current value matches one of the allowed enums.
func (r ReasonConferenceEnded) Validate() error {
	switch r {
	case ReasonEndedViaAPI,
		ReasonParticipantKicked,
		ReasonParticipantLeft,
		ReasonParticipantKickedWithEndOnExit,
		ReasonParticipantLeftWithEndOnExit:
		return nil
	}
	return fmt.Errorf("unknown %T value: %s", r, r)
}

// UnmarshalText implements the TextUnmarshaler interface.
func (r *ReasonConferenceEnded) UnmarshalText(text []byte) (err error) {
	check := ReasonConferenceEnded(text)

	if err := check.Validate(); err != nil {
		return err
	}

	*r = check

	return nil
}

// ConferenceCall is the generic struct for handling common events.
type ConfrenceCall struct {
	ConferenceSID          string `json:"ConferenceSid"`
	FriendlyName           string
	AccountSID             string `json:"AccountSid"`
	SequenceNumber         uint
	Timestamp              TimeRFC1123Z
	StatusCallbackEvent    ConferenceStatus
	CallSID                *string `json:"CallSid,omitempty"`
	Muted                  *bool   `json:",omitempty"`
	Hold                   *bool   `json:",omitempty"`
	Coaching               *bool   `json:",omitempty"`
	EndConferenceOnExit    *bool   `json:",omitempty"`
	StartConferenceOnEnter *bool   `json:",omitempty"`
}

var _ Validator = (*ConfrenceCall)(nil)

// Validate checks that the current value matches one of the allowed enums.
func (c ConfrenceCall) Validate() error {
	switch {
	case !IsValid(c.StatusCallbackEvent):
		return fmt.Errorf("event StatusCallbackEvent invalid got='%s'", c.StatusCallbackEvent)
	}

	return nil
}

// ConfrenceCallEnd is for confrence-end events.
type ConfrenceCallEnd struct {
	ConferenceSID                    string `json:"ConferenceSid"`
	FriendlyName                     string
	AccountSID                       string `json:"AccountSid"`
	SequenceNumber                   uint
	Timestamp                        TimeRFC1123Z
	StatusCallbackEvent              ConferenceStatus
	CallSIDEndingConference          string `json:"CallSidEndingConference"`
	ParticipantLabelEndingConference string
	ReasonConferenceEnded            ReasonConferenceEnded
	Reason                           string
}

var _ Validator = (*ConfrenceCallEnd)(nil)

// Validate checks that the current value matches one of the allowed enums.
func (c ConfrenceCallEnd) Validate() error {
	switch {
	case !IsValid(c.StatusCallbackEvent):
		return fmt.Errorf("event StatusCallbackEvent invalid got='%s'", c.StatusCallbackEvent)

	case c.StatusCallbackEvent == StatusConferenceEnd && c.ReasonConferenceEnded == "":
		return fmt.Errorf("event ReasonConferenceEnded empty")

	case c.StatusCallbackEvent == StatusConferenceEnd && !IsValid(c.ReasonConferenceEnded):
		return fmt.Errorf("event ReasonConferenceEnded invalid got='%s'", c.ReasonConferenceEnded)
	}

	return nil
}

// ConfrenceCallAnnouncement is for confrence-announcement events.
type ConfrenceCallAnnouncement struct {
	ConferenceSID            string `json:"ConferenceSid"`
	FriendlyName             string
	AccountSID               string `json:"AccountSid"`
	SequenceNumber           uint
	Timestamp                TimeRFC1123Z
	StatusCallbackEvent      ConferenceStatus
	CallSID                  *string `json:"CallSid"`
	Muted                    *bool
	Hold                     *bool
	Coaching                 *bool
	EndConferenceOnExit      *bool
	StartConferenceOnEnter   *bool
	ReasonAnnouncementFailed *string
	AnnounceURL              *string `json:"AnnounceUrl"`
}

var _ Validator = (*ConfrenceCallAnnouncement)(nil)

// Validate checks that the current value matches one of the allowed enums.
func (c ConfrenceCallAnnouncement) Validate() error {
	switch {
	case !IsValid(c.StatusCallbackEvent):
		return fmt.Errorf("event StatusCallbackEvent invalid got='%s'", c.StatusCallbackEvent)
	}

	return nil
}
