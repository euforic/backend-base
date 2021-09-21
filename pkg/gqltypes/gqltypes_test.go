package gqltypes

import (
	"bytes"
	"testing"
)

func TestDateTime_UnmarshalGQL(t *testing.T) {
	dt1 := DateTime("2019-10-12T07:20:50.52Z")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		d       *DateTime
		args    args
		wantErr bool
	}{
		{"TestDateTime_UnmarshalGQL", &dt1, args{123}, true},
		{"TestDateTime_UnmarshalGQL", &dt1, args{"2019-10-12T07:20:50.52Z"}, false},
		{"TestDateTime_UnmarshalGQL", &dt1, args{"10:00PM"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("DateTime.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDateTime_MarshalGQL(t *testing.T) {
	dt1 := DateTime("2019-10-12T07:20:50.52Z")
	tests := []struct {
		name  string
		d     DateTime
		wantW string
	}{
		{"TestDateTime_MarshalGQL", dt1, `"2019-10-12T07:20:50.52Z"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.d.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("DateTime.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestHourMinute_UnmarshalGQL(t *testing.T) {
	hm1 := HourMinute("10:00PM")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		h       *HourMinute
		args    args
		wantErr bool
	}{
		{"TestHourMinute_UnmarshalGQL", &hm1, args{123}, true},
		{"TestHourMinute_UnmarshalGQL", &hm1, args{"2019-10-12T07:20:50.52Z"}, true},
		{"TestHourMinute_UnmarshalGQL", &hm1, args{"10:00 AM"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.h.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("HourMinute.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHourMinute_MarshalGQL(t *testing.T) {
	hm1 := HourMinute("10:00AM")
	tests := []struct {
		name  string
		h     HourMinute
		wantW string
	}{
		{"TestHourMinute_MarshalGQL", hm1, `"10:00AM"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.h.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("HourMinute.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
