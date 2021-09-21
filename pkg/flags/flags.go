package flags

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// DeprecatedFlag encapsulates a flag that may have been previously valid but
// is now deprecated. If a DeprecatedFlag is set, an error occurs.
type DeprecatedFlag struct {
	Name string
}

// Set returns an error
func (f *DeprecatedFlag) Set(_ string) error {
	return fmt.Errorf(`flag "-%s" is no longer supported.`, f.Name)
}

// String ...
func (f *DeprecatedFlag) String() string {
	return ""
}

// IgnoredFlag encapsulates a flag that may have been previously valid but is
// now ignored. If an IgnoredFlag is set, a warning is printed and
// operation continues.
type IgnoredFlag struct {
	Name string
}

// IsBoolFlag is defined to allow the flag to be defined without an argument
func (f *IgnoredFlag) IsBoolFlag() bool {
	return true
}

// Set is not supported anymore
func (f *IgnoredFlag) Set(s string) error {
	log.Printf(`flag "-%s" is no longer supported - ignoring.`, f.Name)
	return nil
}

// String ...
func (f *IgnoredFlag) String() string {
	return ""
}

// SetFlagsFromEnv parses all registered flags in the given flagset,
// and if they are not already set it attempts to set their values from
// environment variables. Environment variables take the name of the flag but
// are UPPERCASE, have the given prefix  and any dashes are replaced by
// underscores - for example: some-flag => PREFIX_SOME_FLAG
func SetFlagsFromEnv(prefix string, fs *flag.FlagSet) error {
	var err error
	alreadySet := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		alreadySet[FlagToEnv(prefix, f.Name)] = true
	})
	usedEnvKey := make(map[string]bool)
	fs.VisitAll(func(f *flag.Flag) {
		f.Usage = fmt.Sprintf("%s [%s]", f.Usage, FlagToEnv(prefix, f.Name))
		err = setFlagFromEnv(fs, prefix, f.Name, usedEnvKey, alreadySet, true)
	})

	verifyEnv(prefix, usedEnvKey, alreadySet)

	return err
}

// SetPflagsFromEnv is similar to SetFlagsFromEnv. However, the accepted flagset type is pflag.FlagSet
// and it does not do any logging.
func SetPflagsFromEnv(prefix string, fs *pflag.FlagSet) error {
	var err error
	alreadySet := make(map[string]bool)
	usedEnvKey := make(map[string]bool)
	fs.VisitAll(func(f *pflag.Flag) {
		f.Usage = fmt.Sprintf("%s [%s]", f.Usage, FlagToEnv(prefix, f.Name))
		if f.Changed {
			alreadySet[FlagToEnv(prefix, f.Name)] = true
		}
		if serr := setFlagFromEnv(fs, prefix, f.Name, usedEnvKey, alreadySet, false); serr != nil {
			err = serr
		}
	})
	return err
}

// FlagToEnv converts flag string to upper-case environment variable key string.
func FlagToEnv(prefix, name string) string {
	return prefix + "_" + strings.ToUpper(strings.Replace(name, "-", "_", -1))
}

func verifyEnv(prefix string, usedEnvKey, alreadySet map[string]bool) {
	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		if len(kv) != 2 {
			log.Printf("found invalid env %s", env)
		}
		if usedEnvKey[kv[0]] {
			continue
		}
		if alreadySet[kv[0]] {
			log.Printf("recognized environment variable %s, but unused: shadowed by corresponding flag ", kv[0])
			continue
		}
		if strings.HasPrefix(env, prefix+"_") {
			log.Printf("unrecognized environment variable %s", env)
		}
	}
}

type flagSetter interface {
	Set(fk string, fv string) error
}

func setFlagFromEnv(fs flagSetter, prefix, fname string, usedEnvKey, alreadySet map[string]bool, logging bool) error {
	key := FlagToEnv(prefix, fname)
	if !alreadySet[key] {
		val := os.Getenv(key)
		if val != "" {
			usedEnvKey[key] = true
			if serr := fs.Set(fname, val); serr != nil {
				return fmt.Errorf("invalid value %q for %s: %v", val, key, serr)
			}
			if logging {
				log.Printf("recognized and used environment variable %s=%s", key, val)
			}
		}
	}
	return nil
}

// IsSet checks if flagset has set
func IsSet(fs *flag.FlagSet, name string) bool {
	set := false
	fs.VisitAll(func(f *flag.Flag) {
		if f.Name == name {
			set = true
		}
	})
	return set
}

// IsPSet checks if flagset has set but with pflag
func IsPSet(fs *pflag.FlagSet, name string) bool {
	set := false
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Name == name {
			set = true
		}
	})
	return set
}
