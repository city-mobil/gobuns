package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/pflag"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlices_CfgFile(t *testing.T) {
	cfgPath, err := filepath.Abs("testdata/config.yaml")
	require.NoError(t, err)
	os.Args = append(os.Args, "--config="+cfgPath)

	tests := []struct {
		name string
		want interface{}
		got  interface{}
	}{
		{
			name: "UintSlice",
			want: &[]uint{4, 8},
			got:  UintSlice("uint_slice", []uint{}, ""),
		},
		{
			name: "IntSlice",
			want: &[]int{-1, 8},
			got:  IntSlice("int_slice", []int{}, ""),
		},
		{
			name: "Int64Slice",
			want: &[]int64{6, 4},
			got:  Int64Slice("int64_slice", []int64{}, ""),
		},
		{
			name: "Float32Slice",
			want: &[]float32{9.1, 17.8},
			got:  Float32Slice("float32_slice", []float32{}, ""),
		},
		{
			name: "Float64Slice",
			want: &[]float64{8.6, 32.3},
			got:  Float64Slice("float64_slice", []float64{}, ""),
		},
		{
			name: "StringSlice",
			want: &[]string{"val1", "val2"},
			got:  StringSlice("string_slice", []string{}, ""),
		},
		{
			name: "DurationSlice",
			want: &[]time.Duration{2 * time.Second, 5 * time.Millisecond},
			got:  DurationSlice("duration_slice", []time.Duration{}, ""),
		},
		{
			name: "BoolSlice",
			want: &[]bool{true, false},
			got:  BoolSlice("bool_slice", []bool{}, ""),
		},
		{
			name: "FlagNotFound",
			want: &[]string{"defValue"},
			got:  StringSlice("not_existing_flag", []string{"defValue"}, ""),
		},
		{
			name: "CommaInQuotedElement",
			want: &[]string{"val1", "val2", "val3"},
			got:  StringSlice("string_comma_slice", []string{}, ""),
		},
		{
			name: "CommaInDefaultQuotedElement",
			want: &[]string{"val1", "val2", "val3"},
			got:  StringSlice("not_existing_flag_two", []string{"val1,val2", "val3"}, ""),
		},
		{
			name: "BadTypeInSlice",
			want: &[]int{},
			got:  IntSlice("invalid_int_slice", []int{}, ""),
		},
	}

	err = InitOnce()
	require.NoError(t, err)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.got)
		})
	}
}

func TestRawData(t *testing.T) {
	cfgPath, err := filepath.Abs("testdata/raw_bytes.yaml")
	require.NoError(t, err)

	f := NewFlagSet("default1", pflag.ContinueOnError)
	bytes := f.RawData("raw_bytes", nil, "")

	err = f.Init(&cfgPath)
	require.NoError(t, err)

	assert.Equal(t, []byte(`- name: vasya
  surname: pupkin
- name: agro
  surname: belt
`),
		*bytes)
}
