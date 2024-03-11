package common

import (
	"strconv"
	"time"
)

type Timestamp time.Time

// UnmarshalText defines how encoding/xml unmarshals the object from XML,
// a UNIX timestamp string is converted to int which is used for the Timestamp
// object value
func (t *Timestamp) UnmarshalText(text []byte) error {
	ts, err := strconv.ParseInt(string(text), 10, 0)
	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(ts, 0))

	return nil
}

// MarshalText defines how encoding/xml marshals the object to XML,
// the result is a string of the UNIX timestamp
func (t Timestamp) MarshalText() ([]byte, error) {
	ts := time.Time(t).Unix()
	stamp := strconv.FormatInt(ts, 10)

	return []byte(stamp), nil
}
