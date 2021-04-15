package events

import (
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

end:
	CallSidEndingConference
	ParticipantLabelEndingConference
	ReasonConferenceEnded
	Reason

announcement:
	CallSid
	Muted
	Hold
	Coaching
	EndConferenceOnExit
	StartConferenceOnEnter
	ReasonAnnouncementFailed
	AnnounceUrl


join leave mute hold speaker
	CallSid
	Muted
	Hold
	Coaching
	EndConferenceOnExit
	StartConferenceOnEnter
*/

type ConferenceStatus string

var _ Validator = (*ConferenceStatus)(nil)

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

func (es *ConferenceStatus) UnmarshalText(text []byte) (err error) {
	check := ConferenceStatus(text)

	if err := check.Validate(); err != nil {
		return err
	}

	*es = check

	return nil
}

type ReasonConferenceEnded string

var _ Validator = (*ReasonConferenceEnded)(nil)

const (
	ReasonEndedViaAPI                    ReasonConferenceEnded = "conference-ended-via-api"
	ReasonParticipantKicked              ReasonConferenceEnded = "last-participant-kicked"
	ReasonParticipantLeft                ReasonConferenceEnded = "last-participant-left"
	ReasonParticipantKickedWithEndOnExit ReasonConferenceEnded = "participant-with-end-conference-on-exit-kicked"
	ReasonParticipantLeftWithEndOnExit   ReasonConferenceEnded = "participant-with-end-conference-on-exit-left"
)

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

func (r *ReasonConferenceEnded) UnmarshalText(text []byte) (err error) {
	check := ReasonConferenceEnded(text)

	if err := check.Validate(); err != nil {
		return err
	}

	*r = check

	return nil
}

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

func (c ConfrenceCall) Validate() error {
	switch {
	case !IsValid(c.StatusCallbackEvent):
		return fmt.Errorf("event StatusCallbackEvent invalid got='%s'", c.StatusCallbackEvent)
	}

	return nil
}

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

func (c ConfrenceCallAnnouncement) Validate() error {
	switch {
	case !IsValid(c.StatusCallbackEvent):
		return fmt.Errorf("event StatusCallbackEvent invalid got='%s'", c.StatusCallbackEvent)
	}

	return nil
}
