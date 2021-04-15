package events

import (
	"fmt"
)

type ConferenceEventStatus string

var _ Validator = (*ConferenceEventStatus)(nil)

const (
	StatusConferenceEnd          ConferenceEventStatus = "conference-end"
	StatusConferenceStart        ConferenceEventStatus = "conference-start"
	StatusParticipantLeave       ConferenceEventStatus = "participant-leave"
	StatusParticipantJoin        ConferenceEventStatus = "participant-join"
	StatusParticipantMute        ConferenceEventStatus = "participant-mute"
	StatusParticipantUnmute      ConferenceEventStatus = "participant-unmute"
	StatusParticipantHold        ConferenceEventStatus = "participant-hold"
	StatusParticipantUnhold      ConferenceEventStatus = "participant-unhold"
	StatusParticipantSpeechStart ConferenceEventStatus = "participant-speech-start"
	StatusParticipantSpeechStop  ConferenceEventStatus = "participant-speech-stop"
	StatusAnnouncementEnd        ConferenceEventStatus = "announcement-end"
	StatusAnnouncementFail       ConferenceEventStatus = "announcement-fail"
)

func (es ConferenceEventStatus) Validate() error {
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

func (es *ConferenceEventStatus) UnmarshalText(text []byte) (err error) {
	check := ConferenceEventStatus(text)

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
	StatusCallbackEvent              ConferenceEventStatus
	CallSID                          *string                `json:"CallSid,omitempty"`
	Muted                            *bool                  `json:",omitempty"`
	Hold                             *bool                  `json:",omitempty"`
	Coaching                         *bool                  `json:",omitempty"`
	EndConferenceOnExit              *bool                  `json:",omitempty"`
	StartConferenceOnEnter           *bool                  `json:",omitempty"`
	CallSIDEndingConference          *string                `json:"CallSidEndingConference,omitempty"`
	ParticipantLabelEndingConference *string                `json:",omitempty"`
	ReasonConferenceEnded            *ReasonConferenceEnded `json:",omitempty"`
	Reason                           *string                `json:",omitempty"`
	ReasonAnnouncementFailed         *string                `json:",omitempty"`
	AnnounceURL                      *string                `json:"AnnounceUrl,omitempty"`
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
