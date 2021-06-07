package config

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	defaultFlagSet *FlagSet
)

var (
	configPath *string
)

var (
	sliceTypes = map[string]struct{}{
		"int64Slice":    {},
		"intSlice":      {},
		"stringSlice":   {},
		"uintSlice":     {},
		"durationSlice": {},
		"boolSlice":     {},
		"float32Slice":  {},
		"float64Slice":  {},
	}
)

const (
	typeRawData = "RawData"
)

func init() {
	defaultFlagSet = NewFlagSet("default", pflag.ExitOnError)
	configPath = defaultFlagSet.String("config", "", "Configuration file path.")
}

func marshalBack(x interface{}, configType string) ([]byte, error) {
	if configType == "yaml" || configType == "yml" {
		return yaml.Marshal(x)
	}

	// NOTE(a.petrukhin): default way.
	return json.Marshal(x)
}

func handleRawData(value interface{}, flag *pflag.Flag, configType string) {
	data, err := marshalBack(value, configType)
	if err != nil {
		panic(err)
	}
	_ = flag.Value.Set(string(data))
}

func parseFromConfig(flagSet *FlagSet, configPath string) {
	viper.AutomaticEnv()

	configType := parseConfigType(configPath)
	viper.SetConfigType(configType)
	flagSet.setConfigType(configType)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	flagSet.VisitAll(func(flag *pflag.Flag) {
		x := viper.Get(flag.Name)
		if x == nil {
			return
		}

		switch flag.Value.Type() {
		case typeRawData:
			handleRawData(x, flag, flagSet.configType)
		default:
			_ = flag.Value.Set(getFlagValue(flag))
		}
	})
}

func getFlagValue(flag *pflag.Flag) string {
	if isTypeOfSlice(flag.Value.Type()) {
		return strings.Join(viper.GetStringSlice(flag.Name), ",")
	}

	return viper.GetString(flag.Name)
}

func isTypeOfSlice(flagType string) bool {
	_, ok := sliceTypes[flagType]
	return ok
}

// InitOnce initializes global application configuration.
// All the configuration params MUST be declared before InitOnce is called.
//
// Can be used in such way:
// // var (
// // 	someParam = config.Int("test", 42, "test")
// // )
// //
// // func main() {
// // 	config.InitOnce()
// // }
// It is advised to shutdown application if some error occurred.
func InitOnce() error {
	return defaultFlagSet.Init(configPath, os.Args[1:]...)
}

// FilePath returns path to config file
func FilePath() string {
	return defaultFlagSet.configPath
}

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func Int(name string, defValue int, description string) *int {
	return defaultFlagSet.Int(name, defValue, description)
}

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func (f *FlagSet) Int(name string, defValue int, description string) *int {
	return f.FlagSet.Int(name, defValue, description)
}

// IntSlice defines a []int flag with specified name, default value, and usage string.
// The return value is the address of a []int variable that stores the value of the flag.
func IntSlice(name string, defValue []int, description string) *[]int {
	return defaultFlagSet.IntSlice(name, defValue, description)
}

// IntSlice defines a []int flag with specified name, default value, and usage string.
// The return value is the address of a []int variable that stores the value of the flag.
func (f *FlagSet) IntSlice(name string, defValue []int, description string) *[]int {
	return f.FlagSet.IntSlice(name, defValue, description)
}

// Uint defines a uint flag with specified name, default value, and usage string.
// The return value is the address of a uint  variable that stores the value of the flag.
func Uint(name string, defValue uint, description string) *uint {
	return defaultFlagSet.Uint(name, defValue, description)
}

// Uint defines a uint flag with specified name, default value, and usage string.
// The return value is the address of a uint  variable that stores the value of the flag.
func (f *FlagSet) Uint(name string, defValue uint, description string) *uint {
	return f.FlagSet.Uint(name, defValue, description)
}

// UintSlice defines a []uint flag with specified name, default value, and usage string.
// The return value is the address of a []uint variable that stores the value of the flag.
func UintSlice(name string, defValue []uint, description string) *[]uint {
	return defaultFlagSet.UintSlice(name, defValue, description)
}

// UintSlice defines a []uint flag with specified name, default value, and usage string.
// The return value is the address of a []uint variable that stores the value of the flag.
func (f *FlagSet) UintSlice(name string, defValue []uint, description string) *[]uint {
	return f.FlagSet.UintSlice(name, defValue, description)
}

// Uint8 defines a uint8 flag with specified name, default value, and usage string.
// The return value is the address of a uint8 variable that stores the value of the flag.
func Uint8(name string, defValue uint8, description string) *uint8 {
	return defaultFlagSet.Uint8(name, defValue, description)
}

// Uint8 defines a uint8 flag with specified name, default value, and usage string.
// The return value is the address of a uint8 variable that stores the value of the flag.
func (f *FlagSet) Uint8(name string, defValue uint8, description string) *uint8 {
	return f.FlagSet.Uint8(name, defValue, description)
}

// Uint16 defines a uint flag with specified name, default value, and usage string.
// The return value is the address of a uint  variable that stores the value of the flag.
func Uint16(name string, defValue uint16, description string) *uint16 {
	return defaultFlagSet.Uint16(name, defValue, description)
}

// Uint16 defines a uint flag with specified name, default value, and usage string.
// The return value is the address of a uint  variable that stores the value of the flag.
func (f *FlagSet) Uint16(name string, defValue uint16, description string) *uint16 {
	return f.FlagSet.Uint16(name, defValue, description)
}

// Uint32 defines a uint32 flag with specified name, default value, and usage string.
// The return value is the address of a uint32  variable that stores the value of the flag.
func Uint32(name string, defValue uint32, description string) *uint32 {
	return defaultFlagSet.Uint32(name, defValue, description)
}

// Uint32 defines a uint32 flag with specified name, default value, and usage string.
// The return value is the address of a uint32  variable that stores the value of the flag.
func (f *FlagSet) Uint32(name string, defValue uint32, description string) *uint32 {
	return f.FlagSet.Uint32(name, defValue, description)
}

// Uint64 defines a uint64 flag with specified name, default value, and usage string.
// The return value is the address of a uint64 variable that stores the value of the flag.
func Uint64(name string, defValue uint64, description string) *uint64 {
	return defaultFlagSet.Uint64(name, defValue, description)
}

// Uint64 defines a uint64 flag with specified name, default value, and usage string.
// The return value is the address of a uint64 variable that stores the value of the flag.
func (f *FlagSet) Uint64(name string, defValue uint64, description string) *uint64 {
	return f.FlagSet.Uint64(name, defValue, description)
}

// Int8 defines an int8 flag with specified name, default value, and usage string.
// The return value is the address of an int8 variable that stores the value of the flag.
func Int8(name string, defValue int8, description string) *int8 {
	return defaultFlagSet.Int8(name, defValue, description)
}

// Int8 defines an int8 flag with specified name, default value, and usage string.
// The return value is the address of an int8 variable that stores the value of the flag.
func (f *FlagSet) Int8(name string, defValue int8, description string) *int8 {
	return f.FlagSet.Int8(name, defValue, description)
}

// Int16 defines an int16 flag with specified name, default value, and usage string.
// The return value is the address of an int16 variable that stores the value of the flag.
func Int16(name string, defValue int16, description string) *int16 {
	return defaultFlagSet.Int16(name, defValue, description)
}

// Int16 defines an int16 flag with specified name, default value, and usage string.
// The return value is the address of an int16 variable that stores the value of the flag.
func (f *FlagSet) Int16(name string, defValue int16, description string) *int16 {
	return f.FlagSet.Int16(name, defValue, description)
}

// Int32 defines an int32 flag with specified name, default value, and usage string.
// The return value is the address of an int32 variable that stores the value of the flag.
func Int32(name string, defValue int32, description string) *int32 {
	return defaultFlagSet.Int32(name, defValue, description)
}

// Int32 defines an int32 flag with specified name, default value, and usage string.
// The return value is the address of an int32 variable that stores the value of the flag.
func (f *FlagSet) Int32(name string, defValue int32, description string) *int32 {
	return f.FlagSet.Int32(name, defValue, description)
}

// Int64 defines an int64 flag with specified name, default value, and usage string.
// The return value is the address of an int64 variable that stores the value of the flag.
func Int64(name string, defValue int64, description string) *int64 {
	return defaultFlagSet.Int64(name, defValue, description)
}

// Int64 defines an int64 flag with specified name, default value, and usage string.
// The return value is the address of an int64 variable that stores the value of the flag.
func (f *FlagSet) Int64(name string, defValue int64, description string) *int64 {
	return f.FlagSet.Int64(name, defValue, description)
}

// Int64Slice defines a []int64 flag with specified name, default value, and usage string.
// The return value is the address of a []int64 variable that stores the value of the flag.
func Int64Slice(name string, defValue []int64, description string) *[]int64 {
	return defaultFlagSet.Int64Slice(name, defValue, description)
}

// Int64Slice defines a []int64 flag with specified name, default value, and usage string.
// The return value is the address of a []int64 variable that stores the value of the flag.
func (f *FlagSet) Int64Slice(name string, defValue []int64, description string) *[]int64 {
	return f.FlagSet.Int64Slice(name, defValue, description)
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func String(name, defValue, description string) *string {
	return defaultFlagSet.String(name, defValue, description)
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (f *FlagSet) String(name, defValue, description string) *string {
	return f.FlagSet.String(name, defValue, description)
}

// StringSlice defines a string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
// Compared to StringArray flags, StringSlice flags take comma-separated value as arguments and split them accordingly.
// For example:
//   --ss="v1,v2" --ss="v3"
// will result in
//   []string{"v1", "v2", "v3"}
func StringSlice(name string, defValue []string, description string) *[]string {
	return defaultFlagSet.StringSlice(name, defValue, description)
}

// StringSlice defines a string flag with specified name, default value, and usage string.
// The return value is the address of a []string variable that stores the value of the flag.
// Compared to StringArray flags, StringSlice flags take comma-separated value as arguments and split them accordingly.
// For example:
//   --ss="v1,v2" --ss="v3"
// will result in
//   []string{"v1", "v2", "v3"}
func (f *FlagSet) StringSlice(name string, defValue []string, description string) *[]string {
	return f.FlagSet.StringSlice(name, defValue, description)
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
func Duration(name string, defValue time.Duration, description string) *time.Duration {
	return defaultFlagSet.Duration(name, defValue, description)
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
func (f *FlagSet) Duration(name string, defValue time.Duration, description string) *time.Duration {
	return f.FlagSet.Duration(name, defValue, description)
}

// DurationSlice defines a []time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a []time.Duration variable that stores the value of the flag.
func DurationSlice(name string, defValue []time.Duration, description string) *[]time.Duration {
	return defaultFlagSet.DurationSlice(name, defValue, description)
}

// DurationSlice defines a []time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a []time.Duration variable that stores the value of the flag.
func (f *FlagSet) DurationSlice(name string, defValue []time.Duration, description string) *[]time.Duration {
	return f.FlagSet.DurationSlice(name, defValue, description)
}

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func Bool(name string, defValue bool, description string) *bool {
	return defaultFlagSet.Bool(name, defValue, description)
}

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func (f *FlagSet) Bool(name string, defValue bool, description string) *bool {
	return f.FlagSet.Bool(name, defValue, description)
}

// BoolSlice defines a []bool flag with specified name, default value, and usage string.
// The return value is the address of a []bool variable that stores the value of the flag.
func BoolSlice(name string, defValue []bool, description string) *[]bool {
	return defaultFlagSet.BoolSlice(name, defValue, description)
}

// BoolSlice defines a []bool flag with specified name, default value, and usage string.
// The return value is the address of a []bool variable that stores the value of the flag.
func (f *FlagSet) BoolSlice(name string, defValue []bool, description string) *[]bool {
	return f.FlagSet.BoolSlice(name, defValue, description)
}

// Float32 defines a float32 flag with specified name, default value, and usage string.
// The return value is the address of a float32 variable that stores the value of the flag.
func Float32(name string, defValue float32, description string) *float32 {
	return defaultFlagSet.Float32(name, defValue, description)
}

// Float32 defines a float32 flag with specified name, default value, and usage string.
// The return value is the address of a float32 variable that stores the value of the flag.
func (f *FlagSet) Float32(name string, defValue float32, description string) *float32 {
	return f.FlagSet.Float32(name, defValue, description)
}

// Float32Slice defines a []float32 flag with specified name, default value, and usage string.
// The return value is the address of a []float32 variable that stores the value of the flag.
func Float32Slice(name string, defValue []float32, description string) *[]float32 {
	return defaultFlagSet.Float32Slice(name, defValue, description)
}

// Float32Slice defines a []float32 flag with specified name, default value, and usage string.
// The return value is the address of a []float32 variable that stores the value of the flag.
func (f *FlagSet) Float32Slice(name string, defValue []float32, description string) *[]float32 {
	return f.FlagSet.Float32Slice(name, defValue, description)
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// The return value is the address of a float64 variable that stores the value of the flag.
func Float64(name string, defValue float64, description string) *float64 {
	return defaultFlagSet.Float64(name, defValue, description)
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// The return value is the address of a float64 variable that stores the value of the flag.
func (f *FlagSet) Float64(name string, defValue float64, description string) *float64 {
	return f.FlagSet.Float64(name, defValue, description)
}

// Float64Slice defines a []float64 flag with specified name, default value, and usage string.
// The return value is the address of a []float64 variable that stores the value of the flag.
func Float64Slice(name string, defValue []float64, description string) *[]float64 {
	return defaultFlagSet.Float64Slice(name, defValue, description)
}

// Float64Slice defines a []float64 flag with specified name, default value, and usage string.
// The return value is the address of a []float64 variable that stores the value of the flag.
func (f *FlagSet) Float64Slice(name string, defValue []float64, description string) *[]float64 {
	return f.FlagSet.Float64Slice(name, defValue, description)
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func Var(name string, value pflag.Value, description string) {
	defaultFlagSet.Var(name, value, description)
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (f *FlagSet) Var(name string, value pflag.Value, description string) {
	f.FlagSet.Var(value, name, description)
}

type rawData struct {
	data []byte
}

func (r *rawData) String() string {
	return string(r.data)
}

func (r *rawData) Set(s string) error {
	r.data = []byte(s)
	return nil
}

func (r *rawData) Type() string {
	return typeRawData
}

// RawData defines a raw byte slice for further deserializing on user side with
// specified name, defValue and usage string.
// The return value is the address of a []byte variable that stores the raw-value of the flag.
//
// Raw value is a encoded byte value from configuration file as it is given. If file has yaml format, the data is
// encoded in yaml format, otherwise JSON encoding is used.
func RawData(name string, defValue []byte, description string) *[]byte {
	dt := &rawData{
		data: defValue,
	}

	Var(name, dt, description)
	return &dt.data
}

// RawData defines a raw byte slice for further deserializing on user side with
// specified name, defValue and usage string.
// The return value is the address of a []byte variable that stores the raw-value of the flag.
//
// Raw value is a encoded byte value from configuration file as it is given. If file has yaml format, the data is
// encoded in yaml format, otherwise JSON encoding is used.
func (f *FlagSet) RawData(name string, defValue []byte, description string) *[]byte {
	dt := &rawData{
		data: defValue,
	}

	f.FlagSet.Var(dt, name, description)
	return &dt.data
}

// SubConfigSuffixes gets all the configuration params having one given prefix.
func SubConfigSuffixes(prefix string) []string {
	subConf := viper.GetStringMapStringSlice(prefix)
	suffixes := make([]string, 0, len(subConf))
	for k := range subConf {
		suffixes = append(suffixes, k)
	}
	return suffixes
}

func ReInit(prefix string) {
	defaultFlagSet.ReInit(prefix)
}
