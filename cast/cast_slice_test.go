package cast

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStringSlice(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  []string
		err   error
	}{
		{"InterfaceSlice", []interface{}{1, true}, []string{"1", "true"}, nil},
		{"String", "1 2 3", []string{"1", "2", "3"}, nil},
		{"StringSlice", []string{"val1", "val2"}, []string{"val1", "val2"}, nil},
		{"IntSlice", []int{1, 5}, []string{"1", "5"}, nil},
		{"Float64Slice", []float64{1.3, 5.2}, []string{"1.3", "5.2"}, nil},
		{"ErrorSlice", []error{errors.New("err1"), errors.New("err2")}, []string{"err1", "err2"}, nil},
		{"Bool_ExpectedErr", true, nil, ErrInvalidType},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStringSlice(tt.value)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseBoolSlice(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  []bool
		err   error
	}{
		{"InterfaceSlice", []interface{}{1, "true", 0}, []bool{true, true, false}, nil},
		{"StringSlice", []string{"true", "false", "1"}, []bool{true, false, true}, nil},
		{"IntSlice", []int{1, 5, 0}, []bool{true, true, false}, nil},
		{"Float64Slice", []float64{0, 5.2}, []bool{false, true}, nil},
		{"Int_ExpectedErr", 1, nil, ErrInvalidType},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBoolSlice(tt.value)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseIntSlice(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  []int
		err   error
	}{
		{"InterfaceSlice", []interface{}{1, 2, "3"}, []int{1, 2, 3}, nil},
		{"StringSlice", []string{"1", "2"}, []int{1, 2}, nil},
		{"IntSlice", []int{1, 5, 0}, []int{1, 5, 0}, nil},
		{"BoolSlice_ExpectedErr", []bool{true}, nil, ErrInvalidType},
		{"String_ExpectedErr", "1", nil, ErrInvalidType},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseIntSlice(tt.value)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
