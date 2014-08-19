package govenmo

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Venmo use a time format that's not compatible with Go's default.
// The custom Time type allows parsing their time format.
// It is also designed to be insertable into a Postgres DB (see Scan).
type Time struct {
	time.Time
}

func (venmoTime *Time) Value() (driver.Value, error) {
	if venmoTime == nil {
		return nil, nil
	}
	return venmoTime.Time, nil
}

var timestamptzFormat2 = "2006-01-02 15:04:05.999999999-07"

func (venmoTime *Time) Scan(src interface{}) error {
	stringVal := string(src.([]byte))
	parsedTime, err := time.Parse(timestamptzFormat2, stringVal)
	venmoTime.Time = parsedTime
	return err
}

const VenmoTimeFormat = "2006-01-02T15:04:05"

// UnmarshalJSON([]byte) error

func (venmoTime *Time) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse(VenmoTimeFormat, s)
	if err != nil {
		logger.Println("Couldnt not parse time:", err)
		return err
	}
	venmoTime.Time = t
	return nil
}
