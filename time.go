package notion

import (
	"time"
)

// we need custom Time type to implement custom JSON
// unmarshalling because we need to accept dates in format:
// "2020-12-08T12:00:00Z" and "2020-12-08"
// see https://developers.notion.com/reference/page#date-property-values
// time.Time doesn't unmarshal "2020-12-08"

type Time time.Time

const layout = "2006-01-02"

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format
// or "YYYY-MM-DD" format
func (t *Time) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	st, err := time.Parse(`"`+time.RFC3339+`"`, string(data))
	if err != nil {
		st, err = time.Parse(`"`+layout+`"`, string(data))
	}
	*t = Time(st)
	return err
}

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	// TODO: possibly return "YYYY-MM-DD"
	return time.Time(t).MarshalJSON()
}

// String returns the time in the custom format
func (t *Time) String() string {
	v := time.Time(*t)
	return v.String()
}
