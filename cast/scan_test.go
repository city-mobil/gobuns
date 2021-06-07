package cast

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	var tests = []struct {
		name            string
		tuple           []interface{}
		want            []interface{}
		wantErr         bool
		invalidDestSize bool
	}{
		{
			name: "MultipleTypes",
			tuple: []interface{}{
				"some_string", uint8(1),
			},
			want: []interface{}{
				"some_string", uint8(1),
			},
		},
		{
			name: "UnknownType",
			tuple: []interface{}{
				struct{}{},
			},
			wantErr: true,
		},
		{
			name: "AllKnownTypes",
			tuple: []interface{}{
				"string", uint16(1), uint32(2), uint64(3), uint(4), int8(5), int16(6), int32(7), int64(8),
				int(9), float32(10), float64(11),
			},
			want: []interface{}{
				"string", uint16(1), uint32(2), uint64(3), uint(4), int8(5), int16(6), int32(7), int64(8),
				int(9), float32(10), float64(11),
			},
		},
		{
			name: "InvalidDestinationSize",
			tuple: []interface{}{
				"string",
			},
			wantErr:         true,
			invalidDestSize: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			dest := prepareDestination(tt.tuple)
			if tt.invalidDestSize {
				dest = append(dest, new(int))
			}

			err := Scan(tt.tuple, dest...)
			if tt.wantErr {
				assert.Error(t, err)
			} else { //nolint:gocritic
				if assert.NoError(t, err) {
					for i := 0; i < len(tt.want); i++ {
						g := reflect.Indirect(reflect.ValueOf(dest[i])).Interface()
						assert.Equal(t, tt.want[i], g)
					}
				}
			}
		})
	}
}

func prepareDestination(src []interface{}) []interface{} {
	res := make([]interface{}, 0, len(src))
	for _, v := range src {
		tp := reflect.TypeOf(v)
		res = append(res, reflect.New(tp).Interface())
	}

	return res
}

func BenchmarkScan(b *testing.B) {
	var u uint32
	var v int32
	tuple := []interface{}{
		uint32(1), int32(1),
	}
	for i := 0; i < b.N; i++ {
		err := Scan(tuple, &u, &v)
		if err != nil {
			b.Fatal(err)
		}
	}
}
