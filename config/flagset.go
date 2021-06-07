package config

import (
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/spf13/pflag"
)

const envConfig = "CONFIG"

// A FlagSet represents a set of defined flags.
//
// Current FlagSet is just a small wrapper above pflag.FlagSet
type FlagSet struct {
	*pflag.FlagSet
	once       sync.Once
	configType string
	configPath string
}

// NewFlagSet creates new FlagSet.
func NewFlagSet(name string, errorHandling pflag.ErrorHandling) *FlagSet {
	fs := pflag.NewFlagSet(name, errorHandling)
	return &FlagSet{
		FlagSet: fs,
	}
}

func (f *FlagSet) setConfigType(typ string) {
	f.configType = typ
}

// Init initializes FlagSet.
//
// configPath param is a pointer because it can be overwritten by FlagSet.Parse.
func (f *FlagSet) Init(configPath *string, cmdlineArgs ...string) error {
	var err error
	f.once.Do(func() {
		err = viper.BindPFlags(defaultFlagSet.FlagSet)
		if err != nil {
			return
		}

		err = f.Parse(cmdlineArgs)
		if err != nil {
			return
		}

		if configPath != nil {
			f.configPath = *configPath
		}

		if f.configPath == "" {
			if path, ok := os.LookupEnv(envConfig); ok {
				f.configPath = path
			}
		}

		if f.configPath != "" {
			viper.SetConfigFile(f.configPath)
			// NOTE(a.petrukhin): commandline arguments have more priority than config's ones.
			parseFromConfig(f, f.configPath)
		}
	})

	return err
}

// ReInit reinitializes FlagSet.
func (f *FlagSet) ReInit(prefix string) {
	f.VisitAll(func(flag *pflag.Flag) {
		if !strings.HasPrefix(flag.Name, prefix) {
			return
		}

		x := viper.Get(flag.Name)
		if x == nil {
			return
		}

		switch flag.Value.Type() {
		case typeRawData:
			handleRawData(x, flag, f.configType)
		default:
			_ = flag.Value.Set(getFlagValue(flag))
		}
	})
}
