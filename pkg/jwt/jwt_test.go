package jwt

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Decoder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	type fields struct {
		KeyServerURL string
		DevMode      bool
	}
	type args struct {
		t string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *JSONWebToken
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Decoder{
				KeyServerURL: tt.fields.KeyServerURL,
				DevMode:      tt.fields.DevMode,
			}
			if got := d.Decode(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decoder.Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecoder_Verify(t *testing.T) {
	type fields struct {
		KeyServerURL string
		DevMode      bool
	}
	type args struct {
		t   JSONWebToken
		raw string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   JSONWebToken
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Decoder{
				KeyServerURL: tt.fields.KeyServerURL,
				DevMode:      tt.fields.DevMode,
			}
			if got := d.Verify(tt.args.t, tt.args.raw); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decoder.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *JSONWebToken
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromContext(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	type args struct {
		decoder *Decoder
		next    http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AuthMiddleware(tt.args.decoder, tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
