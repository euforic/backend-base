package flags

import (
	"flag"
	"os"
	"testing"

	"github.com/spf13/pflag"
)

func TestIsSet(t *testing.T) {
	f1 := flag.NewFlagSet("f1", flag.ContinueOnError)
	f1.Bool("silent", true, "das")
	if err := f1.Set("silent", "true"); err != nil {
		t.Error(err)
	}
	f2 := flag.NewFlagSet("f2", flag.ContinueOnError)
	f2.Bool("loud", false, "das")

	type args struct {
		fs   *flag.FlagSet
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "TestIsSet", args: args{fs: f1, name: "silent"}, want: true},
		{name: "TestIsSet", args: args{fs: f2, name: "silent"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSet(tt.args.fs, tt.args.name); got != tt.want {
				t.Errorf("IsSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlagToEnv(t *testing.T) {
	type args struct {
		prefix string
		name   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"TestFlagToEnv1", args{"pre", "database_url"}, "pre_DATABASE_URL"},
		{"TestFlagToEnv2", args{"!@#", "Ten!@#"}, "!@#_TEN!@#"},
		{"TestFlagToEnv3", args{"", ""}, "_"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FlagToEnv(tt.args.prefix, tt.args.name); got != tt.want {
				t.Errorf("FlagToEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetPflagsFromEnv(t *testing.T) {
	os.Setenv("pre_ONE", "true")
	os.Setenv("pre_TWO", "true")
	os.Setenv("no_Te", "false")

	f1 := pflag.NewFlagSet("f1", pflag.ContinueOnError)
	f1.Bool("pre_ONE", false, "das")
	f2 := pflag.NewFlagSet("f2", pflag.ContinueOnError)
	f2.Bool("no_Te", false, "das")

	type args struct {
		prefix string
		fs     *pflag.FlagSet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
		exist   bool
	}{
		{"TestPFlagsFromEnv", args{"pre", f1}, false, "pre_ONE", true},
		{"TestPFlagsFromEnv", args{"pre", f2}, false, "no_Te", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetPflagsFromEnv(tt.args.prefix, tt.args.fs); (err != nil) != tt.wantErr || !(IsPSet(f1, tt.want) == tt.exist) {
				t.Errorf("SetPflagsFromEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetFlagsFromEnv(t *testing.T) {
	os.Setenv("pre_ONE", "true")
	os.Setenv("pre_TWO", "true")
	os.Setenv("no_Te", "false")

	f1 := flag.NewFlagSet("f1", flag.ContinueOnError)
	f1.Bool("pre_ONE", false, "das")
	f2 := flag.NewFlagSet("f2", flag.ContinueOnError)
	f2.Bool("no_Te", false, "das")

	type args struct {
		prefix string
		fs     *flag.FlagSet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
		exist   bool
	}{
		{"TestPFlagsFromEnv", args{"pre", f1}, false, "pre_ONE", true},
		{"TestPFlagsFromEnv", args{"pre", f2}, false, "no_Te", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetFlagsFromEnv(tt.args.prefix, tt.args.fs); (err != nil) != tt.wantErr || !(IsSet(f1, tt.want) == tt.exist) {
				t.Errorf("SetPflagsFromEnv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
