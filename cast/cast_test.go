package cast

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	wantErrKinds = []reflect.Kind{
		reflect.Float32,
		reflect.Float64,
		reflect.String,
		reflect.Struct,
		reflect.Slice,
		reflect.Map,
		reflect.Chan,
		reflect.Complex64,
		reflect.Complex128,
	}
)

func generateTestData() (data []interface{}, testNames []string) {
	data = []interface{}{
		uint(42),
		42,
		uint8(42),
		int8(42),
		uint16(42),
		int16(42),
		uint32(42),
		int32(42),
		uint64(42),
		int64(42),
		float32(42),
		float64(42),
		"test_string",
		struct{}{},
		make([]int, 1),
		make(map[uint8]struct{}),
		make(chan struct{}),
		complex(0, 1),
		complex64(complex(0, 1)),
	}
	testNames = []string{
		"parse_uint",
		"parse_int",
		"parse_uint8",
		"parse_int8",
		"parse_uint16",
		"parse_int16",
		"parse_uint32",
		"parse_int32",
		"parse_uint64",
		"parse_int64",
		"parse_float32",
		"parse_float64",
		"parse_string",
		"parse_struct",
		"parse_slice",
		"parse_map",
		"parse_chan",
		"parse_complex128",
		"parse_complex64",
	}
	return data, testNames
}

func needError(curr reflect.Kind, exclude ...reflect.Kind) bool {
	result := false
	for _, v := range wantErrKinds {
		if v == curr {
			result = true
			break
		}
	}
	for _, v := range exclude {
		if v == curr {
			result = false
			break
		}
	}
	return result
}

func TestParseInt(t *testing.T) {
	type testCase struct {
		testName string
		data     interface{}
		expected int
		wantErr  bool
	}
	var testData []testCase

	data, testNames := generateTestData()
	for i := 0; i < len(data); i++ {
		typ := reflect.TypeOf(data[i])
		cas := testCase{
			testName: testNames[i],
			data:     data[i],
			expected: 42,
		}
		if needError(typ.Kind()) {
			cas.wantErr = true
			cas.expected = 0
		}
		testData = append(testData, cas)
	}

	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			got, err := ParseInt(v.data)
			if err != nil && !v.wantErr {
				t.Errorf("got err %v, but expected <nil>", err)
			}
			if err == nil && v.wantErr {
				t.Error("got err <nil>, but expected non-nil")
			}
			if got != v.expected {
				t.Errorf("got %d, expected %d", got, v.expected)
			}
		})
	}
}

func TestParseUint(t *testing.T) {
	// NOTE(a.petrukhin): we do not test any overflows.
	type testCase struct {
		name  string
		value interface{}
		want  uint
		err   error
	}

	tests := []testCase{
		{"Int", 2, 2, nil},
		{"Int8", int8(2), 2, nil},
		{"Int16", int16(2), 2, nil},
		{"Int32", int32(2), 2, nil},
		{"Int64", int64(2), 2, nil},
		{"Uint", uint(2), 2, nil},
		{"Uint8", uint8(2), 2, nil},
		{"Uint16", uint16(2), 2, nil},
		{"Uint32", uint32(2), 2, nil},
		{"Uint64", uint64(2), 2, nil},
		{"String", "8", 8, nil},
		{"NegativeInt_ExpectedErr", -14, 0, ErrNegativeToUnsigned},
		{"Bool_ExpectedError", false, 0, ErrInvalidType},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUint(tt.value)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseUint8(t *testing.T) {
	// NOTE(a.petrukhin): we do not test any overflows.
	type testCase struct {
		name  string
		value interface{}
		want  uint8
		err   error
	}

	tests := []testCase{
		{"Int", 2, 2, nil},
		{"Int8", int8(2), 2, nil},
		{"Int16", int16(2), 2, nil},
		{"Int32", int32(2), 2, nil},
		{"Int64", int64(2), 2, nil},
		{"Uint", uint(2), 2, nil},
		{"Uint8", uint8(2), 2, nil},
		{"Uint16", uint16(2), 2, nil},
		{"Uint32", uint32(2), 2, nil},
		{"Uint64", uint64(2), 2, nil},
		{"String", "8", 8, nil},
		{"NegativeInt_ExpectedErr", -14, 0, ErrNegativeToUnsigned},
		{"Bool_ExpectedError", false, 0, ErrInvalidType},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUint8(tt.value)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseInt8(t *testing.T) {
	// NOTE(a.petrukhin): we do not test any overflows.
	type testCase struct {
		testName string
		data     interface{}
		expected int8
		wantErr  bool
	}
	var testData []testCase

	data, testNames := generateTestData()
	for i := 0; i < len(data); i++ {
		typ := reflect.TypeOf(data[i])
		cas := testCase{
			testName: testNames[i],
			data:     data[i],
			expected: 42,
		}
		if needError(typ.Kind()) {
			cas.wantErr = true
			cas.expected = 0
		}
		testData = append(testData, cas)
	}

	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			got, err := ParseInt8(v.data)
			if err != nil && !v.wantErr {
				t.Errorf("got err %v, but expected <nil>", err)
			}
			if err == nil && v.wantErr {
				t.Error("got err <nil>, but expected non-nil")
			}
			if got != v.expected {
				t.Errorf("got %d, expected %d", got, v.expected)
			}
		})
	}
}

func TestParseInt16(t *testing.T) {
	// NOTE(a.petrukhin): we do not test any overflows.
	type testCase struct {
		testName string
		data     interface{}
		expected int16
		wantErr  bool
	}
	var testData []testCase

	data, testNames := generateTestData()
	for i := 0; i < len(data); i++ {
		typ := reflect.TypeOf(data[i])
		cas := testCase{
			testName: testNames[i],
			data:     data[i],
			expected: 42,
		}
		if needError(typ.Kind()) {
			cas.wantErr = true
			cas.expected = 0
		}
		testData = append(testData, cas)
	}

	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			got, err := ParseInt16(v.data)
			if err != nil && !v.wantErr {
				t.Errorf("got err %v, but expected <nil>", err)
			}
			if err == nil && v.wantErr {
				t.Error("got err <nil>, but expected non-nil")
			}
			if got != v.expected {
				t.Errorf("got %d, expected %d", got, v.expected)
			}
		})
	}
}

func TestParseUint16(t *testing.T) {
	// NOTE(a.petrukhin): we do not test any overflows.
	type testCase struct {
		name  string
		value interface{}
		want  uint16
		err   error
	}

	tests := []testCase{
		{"Int", 2, 2, nil},
		{"Int8", int8(2), 2, nil},
		{"Int16", int16(2), 2, nil},
		{"Int32", int32(2), 2, nil},
		{"Int64", int64(2), 2, nil},
		{"Uint", uint(2), 2, nil},
		{"Uint8", uint8(2), 2, nil},
		{"Uint16", uint16(2), 2, nil},
		{"Uint32", uint32(2), 2, nil},
		{"Uint64", uint64(2), 2, nil},
		{"String", "8", 8, nil},
		{"NegativeInt_ExpectedErr", -14, 0, ErrNegativeToUnsigned},
		{"Bool_ExpectedError", false, 0, ErrInvalidType},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUint16(tt.value)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseInt32(t *testing.T) {
	// NOTE(a.petrukhin): we do not test any overflows.
	type testCase struct {
		testName string
		data     interface{}
		expected int32
		wantErr  bool
	}
	var testData []testCase

	data, testNames := generateTestData()
	for i := 0; i < len(data); i++ {
		typ := reflect.TypeOf(data[i])
		cas := testCase{
			testName: testNames[i],
			data:     data[i],
			expected: 42,
		}
		if needError(typ.Kind()) {
			cas.wantErr = true
			cas.expected = 0
		}
		testData = append(testData, cas)
	}

	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			got, err := ParseInt32(v.data)
			if err != nil && !v.wantErr {
				t.Errorf("got err %v, but expected <nil>", err)
			}
			if err == nil && v.wantErr {
				t.Error("got err <nil>, but expected non-nil")
			}
			if got != v.expected {
				t.Errorf("got %d, expected %d", got, v.expected)
			}
		})
	}
}

func TestParseUint32(t *testing.T) {
	// NOTE(a.petrukhin): we do not test any overflows.
	type testCase struct {
		name  string
		value interface{}
		want  uint32
		err   error
	}

	tests := []testCase{
		{"Int", 2, 2, nil},
		{"Int8", int8(2), 2, nil},
		{"Int16", int16(2), 2, nil},
		{"Int32", int32(2), 2, nil},
		{"Int64", int64(2), 2, nil},
		{"Uint", uint(2), 2, nil},
		{"Uint8", uint8(2), 2, nil},
		{"Uint16", uint16(2), 2, nil},
		{"Uint32", uint32(2), 2, nil},
		{"Uint64", uint64(2), 2, nil},
		{"String", "8", 8, nil},
		{"NegativeInt_ExpectedErr", -14, 0, ErrNegativeToUnsigned},
		{"Bool_ExpectedError", false, 0, ErrInvalidType},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUint32(tt.value)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseUint64(t *testing.T) {
	// NOTE(a.petrukhin): we do not test any overflows.
	type testCase struct {
		name  string
		value interface{}
		want  uint64
		err   error
	}

	tests := []testCase{
		{"Int", 2, 2, nil},
		{"Int8", int8(2), 2, nil},
		{"Int16", int16(2), 2, nil},
		{"Int32", int32(2), 2, nil},
		{"Int64", int64(2), 2, nil},
		{"Uint", uint(2), 2, nil},
		{"Uint8", uint8(2), 2, nil},
		{"Uint16", uint16(2), 2, nil},
		{"Uint32", uint32(2), 2, nil},
		{"Uint64", uint64(2), 2, nil},
		{"String", "8", 8, nil},
		{"NegativeInt_ExpectedErr", -14, 0, ErrNegativeToUnsigned},
		{"Bool_ExpectedError", false, 0, ErrInvalidType},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUint64(tt.value)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseInt64(t *testing.T) {
	// NOTE(a.petrukhin): we do not test any overflows.
	type testCase struct {
		testName string
		data     interface{}
		expected int64
		wantErr  bool
	}
	var testData []testCase

	data, testNames := generateTestData()
	for i := 0; i < len(data); i++ {
		typ := reflect.TypeOf(data[i])
		cas := testCase{
			testName: testNames[i],
			data:     data[i],
			expected: 42,
		}
		if needError(typ.Kind()) {
			cas.wantErr = true
			cas.expected = 0
		}
		testData = append(testData, cas)
	}

	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			got, err := ParseInt64(v.data)
			if err != nil && !v.wantErr {
				t.Errorf("got err %v, but expected <nil>", err)
			}
			if err == nil && v.wantErr {
				t.Error("got err <nil>, but expected non-nil")
			}
			if got != v.expected {
				t.Errorf("got %d, expected %d", got, v.expected)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string
		err   error
	}{
		{"Int", 2, "2", nil},
		{"Int8", int8(2), "2", nil},
		{"Int16", int16(2), "2", nil},
		{"Int32", int32(2), "2", nil},
		{"Int64", int64(2), "2", nil},
		{"Uint", uint(2), "2", nil},
		{"Uint8", uint8(2), "2", nil},
		{"Uint16", uint16(2), "2", nil},
		{"Uint32", uint32(2), "2", nil},
		{"Uint64", uint64(2), "2", nil},
		{"Float32", float32(2.5), "2.5", nil},
		{"Float64", 2.6, "2.6", nil},
		{"String", "2", "2", nil},
		{"Bool", true, "true", nil},
		{"Nil", nil, "", nil},
		{"SliceByte", []byte(`test`), "test", nil},
		{"Error", errors.New(`test`), "test", nil},
		{"SliceInt_ExpectedError", []int{1}, "", ErrInvalidType},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseString(tt.value)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
		err   error
	}{
		{"Int", 2, true, nil},
		{"Int8", int8(-2), true, nil},
		{"Int16", int16(0), false, nil},
		{"Int32", int32(2), true, nil},
		{"Int64", int64(2), true, nil},
		{"Uint", uint(0), false, nil},
		{"Uint8", uint8(2), true, nil},
		{"Uint16", uint16(2), true, nil},
		{"Uint32", uint32(2), true, nil},
		{"Uint64", uint64(0), false, nil},
		{"Float32", float32(2.5), true, nil},
		{"Float64", 2.6, true, nil},
		{"String", "true", true, nil},
		{"Bool", true, true, nil},
		{"Nil", nil, false, nil},
		{"SliceInt_ExpectedError", []int{1}, false, ErrInvalidType},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBool(tt.value)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseFloat32(t *testing.T) {
	type testCase struct {
		testName string
		data     interface{}
		expected float32
		wantErr  bool
	}
	var testData []testCase

	data, testNames := generateTestData()
	for i := 0; i < len(data); i++ {
		typ := reflect.TypeOf(data[i])
		cas := testCase{
			testName: testNames[i],
			data:     data[i],
			expected: 42,
		}
		if needError(typ.Kind(), reflect.Float32, reflect.Float64) {
			cas.wantErr = true
			cas.expected = 0
		}
		testData = append(testData, cas)
	}
	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			got, err := ParseFloat32(v.data)
			if err != nil && !v.wantErr {
				t.Errorf("got err %v, but expected <nil>", err)
			}
			if err == nil && v.wantErr {
				t.Error("got err <nil>, but expected non-nil")
			}
			if got != v.expected {
				t.Errorf("got %v, expected %v", got, v.expected)
			}
		})
	}
}

func TestParseFloat64(t *testing.T) {
	type testCase struct {
		testName string
		data     interface{}
		expected float64
		wantErr  bool
	}
	var testData []testCase

	data, testNames := generateTestData()
	for i := 0; i < len(data); i++ {
		typ := reflect.TypeOf(data[i])
		cas := testCase{
			testName: testNames[i],
			data:     data[i],
			expected: 42,
		}
		if needError(typ.Kind(), reflect.Float64, reflect.Float32) {
			cas.wantErr = true
			cas.expected = 0
		}
		testData = append(testData, cas)
	}
	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			got, err := ParseFloat64(v.data)
			if err != nil && !v.wantErr {
				t.Errorf("got err %v, but expected <nil>", err)
			}
			if err == nil && v.wantErr {
				t.Error("got err <nil>, but expected non-nil")
			}
			if got != v.expected {
				t.Errorf("got %v, expected %v", got, v.expected)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  time.Duration
		err   error
	}{
		{"Int", 2, 2 * time.Nanosecond, nil},
		{"Int8", int8(2), 2 * time.Nanosecond, nil},
		{"Int16", int16(2), 2 * time.Nanosecond, nil},
		{"Int32", int32(2), 2 * time.Nanosecond, nil},
		{"Int64", int64(2), 2 * time.Nanosecond, nil},
		{"Uint", uint(2), 2 * time.Nanosecond, nil},
		{"Uint8", uint8(2), 2 * time.Nanosecond, nil},
		{"Uint16", uint16(2), 2 * time.Nanosecond, nil},
		{"Uint32", uint32(2), 2 * time.Nanosecond, nil},
		{"Uint64", uint64(2), 2 * time.Nanosecond, nil},
		{"Float32", float32(8.5), 8 * time.Nanosecond, nil},
		{"Float64", 7.5, 7 * time.Nanosecond, nil},
		{"String", "8ms", 8 * time.Millisecond, nil},
		{"NegativeInt", -14, -14 * time.Nanosecond, nil},
		{"Bool_ExpectedError", false, 0, ErrInvalidType},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.value)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkParseInt64(b *testing.B) {
	v := interface{}(uint64(42))
	for i := 0; i < b.N; i++ {
		_, _ = ParseInt64(v)
	}
}
