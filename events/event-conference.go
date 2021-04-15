package events

import (
	"fmt"
)

type EventStatus string

var _ Validator = (*EventStatus)(nil)

const (
	StatusConferenceEnd          EventStatus = "conference-end"
	StatusConferenceStart        EventStatus = "conference-start"
	StatusParticipantLeave       EventStatus = "participant-leave"
	StatusParticipantJoin        EventStatus = "participant-join"
	StatusParticipantMute        EventStatus = "participant-mute"
	StatusParticipantUnmute      EventStatus = "participant-unmute"
	StatusParticipantHold        EventStatus = "participant-hold"
	StatusParticipantUnhold      EventStatus = "participant-unhold"
	StatusParticipantSpeechStart EventStatus = "participant-speech-start"
	StatusParticipantSpeechStop  EventStatus = "participant-speech-stop"
	StatusAnnouncementEnd        EventStatus = "announcement-end"
	StatusAnnouncementFail       EventStatus = "announcement-fail"
)

func (es EventStatus) Validate() error {
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

func (es *EventStatus) UnmarshalText(text []byte) (err error) {
	check := EventStatus(text)

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

type ConfrenceCallEvent struct {
	ConferenceSID                    string `json:"ConferenceSid"`
	FriendlyName                     string
	AccountSID                       string `json:"AccountSid"`
	SequenceNumber                   uint
	Timestamp                        TimeRFC1123Z
	StatusCallbackEvent              EventStatus
	CallSID                          *string `json:"CallSid,omitempty"`
	Muted                            *bool `json:",omitempty"`
	Hold                             *bool `json:",omitempty"`
	Coaching                         *bool `json:",omitempty"`
	EndConferenceOnExit              *bool `json:",omitempty"`
	StartConferenceOnEnter           *bool `json:",omitempty"`
	CallSIDEndingConference          *string `json:"CallSidEndingConference,omitempty"`
	ParticipantLabelEndingConference *string `json:",omitempty"`
	ReasonConferenceEnded            *ReasonConferenceEnded `json:",omitempty"`
	Reason                           *string `json:",omitempty"`
	ReasonAnnouncementFailed         *string `json:",omitempty"`
	AnnounceURL                      *string `json:"AnnounceUrl,omitempty"`
}

func (c ConfrenceCallEvent) Validate() error {
	switch {
	case !IsValid(c.StatusCallbackEvent):
		return fmt.Errorf("event StatusCallbackEvent invalid got='%s'", c.StatusCallbackEvent)

	case c.StatusCallbackEvent == StatusConferenceEnd && c.ReasonConferenceEnded == nil:
		return fmt.Errorf("event ReasonConferenceEnded empty")
	case c.StatusCallbackEvent == StatusConferenceEnd && c.ReasonConferenceEnded != nil && !IsValid(c.ReasonConferenceEnded):
		return fmt.Errorf("event ReasonConferenceEnded invalid got='%s'", *c.ReasonConferenceEnded)
	}

	return nil
}

type RecordingStatus struct {
}
