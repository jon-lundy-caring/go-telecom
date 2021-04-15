package events_test

import (
	"encoding/json"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jon-lundy-caring/go-telecom/events"

	"github.com/matryer/is"
)

type testDecodeCase struct {
	container func() events.Validator
	form      string
	json      []byte
	err       string
}

func TestConfrenceCallEvents(t *testing.T) {
	tests := []testDecodeCase{
		{
			container: func() events.Validator { return &events.ConfrenceCall{} },
			form:      "ConferenceSid=SID1234&FriendlyName=call-name&AccountSid=AID1234&muted=true&StatusCallbackEvent=participant-mute&Timestamp=Mon, 02 Jan 2006 15:04:05 -0700&CallSid=CA1234",
			json:      []byte(`{"ConferenceSid":"SID1234","FriendlyName":"call-name","AccountSid":"AID1234","SequenceNumber":0,"Timestamp":"Mon, 02 Jan 2006 15:04:05 -0700","StatusCallbackEvent":"participant-mute","CallSid":"CA1234","Muted":true}`),
		},
		{
			container: func() events.Validator { return &events.ConfrenceCallEnd{} },
			form:      "ConferenceSid=SID1234&FriendlyName=call-name&AccountSid=AID1234&StatusCallbackEvent=conference-end&Timestamp=Mon, 02 Jan 2006 15:04:05 -0700&CallSidEndingConference=SID1234&ParticipantLabelEndingConference=PID1234&ReasonConferenceEnded=conference-ended-via-api&Reason=ended+by+host",
			json:      []byte(`{"ConferenceSid":"SID1234","FriendlyName":"call-name","AccountSid":"AID1234","SequenceNumber":0,"Timestamp":"Mon, 02 Jan 2006 15:04:05 -0700","StatusCallbackEvent":"conference-end","CallSidEndingConference":"SID1234","ParticipantLabelEndingConference":"PID1234","ReasonConferenceEnded":"conference-ended-via-api","Reason":"ended by host"}`),
		},
		{
			container: func() events.Validator { return &events.ConfrenceCallAnnouncement{} },
			form:      "ConferenceSid=SID1234&FriendlyName=call-name&AccountSid=AID1234&Timestamp=Mon, 02 Jan 2006 15:04:05 -0700&StatusCallbackEvent=announcement-fail&CallSid=CA1234&muted=true&Hold=false&Coaching=false&EndConferenceOnExit=true&StartConferenceOnEnter=false&ReasonAnnouncementFailed=timeout&AnnounceUrl=http://some.url/file.mp4",
			json:      []byte(`{"ConferenceSid":"SID1234","FriendlyName":"call-name","AccountSid":"AID1234","SequenceNumber":0,"Timestamp":"Mon, 02 Jan 2006 15:04:05 -0700","StatusCallbackEvent":"announcement-fail","CallSid":"CA1234","Muted":true,"Hold":false,"Coaching":false,"EndConferenceOnExit":true,"StartConferenceOnEnter":false,"ReasonAnnouncementFailed":"timeout","AnnounceUrl":"http://some.url/file.mp4"}`),
		},
		{
			container: func() events.Validator { return &events.ConfrenceCallEnd{} },
			form:      "ConferenceSid=SID1234&FriendlyName=call-name&AccountSid=AID1234&StatusCallbackEvent=conference-endXX&Timestamp=Mon, 02 Jan 2006 15:04:05 -0700&CallSidEndingConference=SID1234&ParticipantLabelEndingConference=PID1234&ReasonConferenceEnded=conference-ended-via-api&Reason=ended+by+host",
			json:      []byte(`{"ConferenceSid":"SID1234","FriendlyName":"call-name","AccountSid":"AID1234","SequenceNumber":0,"Timestamp":"Mon, 02 Jan 2006 15:04:05 -0700","StatusCallbackEvent":"conference-endXX","CallSid":"CA1234","Muted":true,"CallSidEndingConference":"SID1234","ParticipantLabelEndingConference":"PID1234","ReasonConferenceEnded":"conference-ended-via-api","Reason":"ended by host"}`),
			err:       `unknown events.ConferenceStatus value: conference-endXX`,
		},
		{
			container: func() events.Validator { return &events.ConfrenceCallEnd{} },
			form:      "ConferenceSid=SID1234&FriendlyName=call-name&AccountSid=AID1234&StatusCallbackEvent=conference-end&Timestamp=Mon, 02 Jan 2006 15:04:05 -0700&CallSidEndingConference=SID1234&ParticipantLabelEndingConference=PID1234&ReasonConferenceEnded=conference-ended-via-apiXX&Reason=ended+by+host",
			json:      []byte(`{"ConferenceSid":"SID1234","FriendlyName":"call-name","AccountSid":"AID1234","SequenceNumber":0,"Timestamp":"Mon, 02 Jan 2006 15:04:05 -0700","StatusCallbackEvent":"conference-end","CallSid":"CA1234","Muted":true,"CallSidEndingConference":"SID1234","ParticipantLabelEndingConference":"PID1234","ReasonConferenceEnded":"conference-ended-via-apiXX","Reason":"ended by host"}`),
			err:       `unknown events.ReasonConferenceEnded value: conference-ended-via-apiXX`,
		},
		{
			container: func() events.Validator { return &events.ConfrenceCallEnd{} },
			form:      "ConferenceSid=SID1234&FriendlyName=call-name&AccountSid=AID1234&StatusCallbackEvent=conference-end&Timestamp=Mon, 02 Jan 2006 15:04:05 -0700&CallSidEndingConference=SID1234&ParticipantLabelEndingConference=PID1234&Reason=ended+by+host",
			json:      []byte(`{"ConferenceSid":"SID1234","FriendlyName":"call-name","AccountSid":"AID1234","SequenceNumber":0,"Timestamp":"Mon, 02 Jan 2006 15:04:05 -0700","StatusCallbackEvent":"conference-end","CallSid":"CA1234","Muted":true,"CallSidEndingConference":"SID1234","ParticipantLabelEndingConference":"PID1234","Reason":"ended by host"}`),
			err:       `event ReasonConferenceEnded empty`,
		},
		{
			container: func() events.Validator { return &events.Recording{} },
			form:      "RecordingSid=RID1234&ConferenceSid=SID1234&AccountSid=AID1234&RecordingUrl=http://example.com/file.mp4&RecordingStatus=completed&RecordingDuration=876&RecordingChannels=1&RecordingStartTime=Mon, 02 Jan 2006 15:04:05 -0700&RecordingSource=Conference",
			json:      []byte(`{"AccountSid":"AID1234","ConferenceSid":"SID1234","RecordingSid":"RID1234","RecordingUrl":"http://example.com/file.mp4","RecordingStatus":"completed","RecordingDuration":876,"RecordingChannels":1,"RecordingStartTime":"Mon, 02 Jan 2006 15:04:05 -0700","RecordingSource":"Conference"}`),
		},
		{
			container: func() events.Validator { return &events.Recording{} },
			form:      "RecordingSid=RID1234&ConferenceSid=SID1234&AccountSid=AID1234&RecordingUrl=http://example.com/file.mp4&RecordingStatus=completedXX&RecordingDuration=876&RecordingChannels=1&RecordingStartTime=Mon, 02 Jan 2006 15:04:05 -0700&RecordingSource=Conference",
			json:      []byte(`{"AccountSid":"AID1234","ConferenceSid":"SID1234","RecordingSid":"RID1234","RecordingUrl":"http://example.com/file.mp4","RecordingStatus":"completedXX","RecordingDuration":876,"RecordingChannels":1,"RecordingStartTime":"Mon, 02 Jan 2006 15:04:05 -0700","RecordingSource":"Conference"}`),
			err:       `unknown events.RecordingStatus value: completedXX`,

		},
	}

	testDecodeFromQueryAndJSON(t, tests)
}

func testDecodeFromQueryAndJSON(t *testing.T, tests []testDecodeCase) {
	t.Helper()

	is := is.New(t)

	for i, tt := range tests {
		t.Logf("DecodeFromQuery Test %d", i)
		o := tt.container()

		// convert form data into url.Values
		v, err := url.ParseQuery(tt.form)
		is.NoErr(err)

		// Testing Decode from Query
		err = events.DecodeFromQuery(v, o)

		if tt.err != "" {
			// Tests for error cases
			is.True(err != nil)
			is.True(!events.IsValid(o))

			t.Log(err)
			is.True(strings.Contains(err.Error(), tt.err))

			// Testing Decode from Json matches
			jo := tt.container()

			err := events.DecodeFromJson(tt.json, jo)
			t.Log(err)
			is.True(strings.Contains(err.Error(), tt.err))
		} else {
			// Tests for non-error cases
			is.NoErr(err)
			is.True(events.IsValid(o))

			j, err := json.Marshal(o)
			is.NoErr(err)
			is.Equal(string(j), string(tt.json))

			// Testing Decode from Json matches
			jo := tt.container()

			err = events.DecodeFromJson(tt.json, jo)
			is.NoErr(err)

			j, err = json.Marshal(jo)
			is.NoErr(err)
			is.Equal(string(j), string(tt.json))
		}
	}
}

func TestTimeRFC1123Z(t *testing.T) {
	res := events.NewTimeRFC1123Z("Mon, 02 Jan 2006 15:04:05 -0700")
	check, _ := time.Parse(time.RFC1123Z, "Mon, 02 Jan 2006 15:04:05 -0700")
	is := is.New(t)

	is.Equal(res.Time(), check)
}
