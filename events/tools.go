package events

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/schema"
)

func DecodeFromQuery(values url.Values, obj Validator) error {
	decoder := schema.NewDecoder()
	decoder.ZeroEmpty(true)
	err := decoder.Decode(obj, values)
	if err != nil {
		return err
	}
	if err = obj.Validate(); err != nil {
		return fmt.Errorf("decoded error: %w", err)
	}
	return nil
}

func DecodeFromJson(b []byte, obj Validator) error {
	err := json.Unmarshal(b, obj)
	if err != nil {
		return err
	}
	if err = obj.Validate(); err != nil {
		return fmt.Errorf("decoded error: %w", err)
	}
	return nil
}

type Validator interface {
	Validate() error
}

func IsValid(v Validator) bool {
	if err := v.Validate(); err != nil {
		return false
	}

	return true
}

type TimeRFC1123Z time.Time

func NewTimeRFC1123Z(s string) TimeRFC1123Z {
	t := TimeRFC1123Z{}
	_ = t.UnmarshalText([]byte(s))
	return t
}

func (r TimeRFC1123Z) Time() time.Time {
	return time.Time(r)
}

func (r TimeRFC1123Z) String() string {
	t := time.Time(r)
	return t.Format(time.RFC1123Z)
}

func (r TimeRFC1123Z) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *TimeRFC1123Z) UnmarshalText(text []byte) error {
	s := string(text)

	t, err := time.Parse(time.RFC1123Z, s)

	*r = TimeRFC1123Z(t)

	return err
}
