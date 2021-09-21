package gqltypes

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// DateTime ...
type DateTime string

// HourMinute  ...
type HourMinute string

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (d *DateTime) UnmarshalGQL(v interface{}) error {
	datetime, ok := v.(string)
	if !ok {
		return fmt.Errorf("DateTime must be a string")
	}

	_, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return fmt.Errorf("DateTime must match RFC3339")
	}

	*d = DateTime(datetime)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (d DateTime) MarshalGQL(w io.Writer) {
	if _, err := w.Write([]byte(strconv.Quote(string(d)))); err != nil {
		panic(err)
	}
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (h *HourMinute) UnmarshalGQL(v interface{}) error {
	hourMinute, ok := v.(string)
	if !ok {
		return fmt.Errorf("HourMinute must be a string")
	}

	_, err := time.Parse(time.Kitchen, hourMinute)
	if err != nil {
		return errors.New("time must be in format of HH:mmAM(PM)")
	}

	*h = HourMinute(hourMinute)
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (h HourMinute) MarshalGQL(w io.Writer) {
	if _, err := w.Write([]byte(strconv.Quote(string(h)))); err != nil {
		os.Stderr.WriteString(err.Error())
	}
}
