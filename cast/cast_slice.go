//nolint:dupl
package cast

import (
	"reflect"
	"strings"
)

func ParseIntSlice(i interface{}) ([]int, error) {
	if i == nil {
		return nil, nil
	}

	val, ok := i.([]int)
	if ok {
		return val, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]int, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := ParseInt(s.Index(j).Interface())
			if err != nil {
				return nil, err
			}
			a[j] = val
		}
		return a, nil
	default:
		return nil, ErrInvalidType
	}
}

func ParseBoolSlice(i interface{}) ([]bool, error) {
	if i == nil {
		return nil, ErrInvalidType
	}

	val, ok := i.([]bool)
	if ok {
		return val, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]bool, s.Len())
		for j := 0; j < s.Len(); j++ {
			val, err := ParseBool(s.Index(j).Interface())
			if err != nil {
				return nil, err
			}
			a[j] = val
		}
		return a, nil
	default:
		return nil, ErrInvalidType
	}
}

func ParseStringSlice(i interface{}) ([]string, error) {
	switch val := i.(type) {
	case []interface{}:
		resp := make([]string, 0, len(val))
		for _, v := range val {
			value, err := ParseString(v)
			if err != nil {
				return nil, err
			}
			resp = append(resp, value)
		}
		return resp, nil
	case string:
		return strings.Fields(val), nil
	case []string:
		return val, nil
	case []int:
		resp := make([]string, 0, len(val))
		for _, v := range val {
			value, err := ParseString(v)
			if err != nil {
				return nil, err
			}
			resp = append(resp, value)
		}
		return resp, nil
	case []int8:
		resp := make([]string, 0, len(val))
		for _, v := range val {
			value, err := ParseString(v)
			if err != nil {
				return nil, err
			}
			resp = append(resp, value)
		}
		return resp, nil
	case []int16:
		resp := make([]string, 0, len(val))
		for _, v := range val {
			value, err := ParseString(v)
			if err != nil {
				return nil, err
			}
			resp = append(resp, value)
		}
		return resp, nil
	case []int32:
		resp := make([]string, 0, len(val))
		for _, v := range val {
			value, err := ParseString(v)
			if err != nil {
				return nil, err
			}
			resp = append(resp, value)
		}
		return resp, nil
	case []int64:
		resp := make([]string, 0, len(val))
		for _, v := range val {
			value, err := ParseString(v)
			if err != nil {
				return nil, err
			}
			resp = append(resp, value)
		}
		return resp, nil
	case []float32:
		resp := make([]string, 0, len(val))
		for _, v := range val {
			value, err := ParseString(v)
			if err != nil {
				return nil, err
			}
			resp = append(resp, value)
		}
		return resp, nil
	case []float64:
		resp := make([]string, 0, len(val))
		for _, v := range val {
			value, err := ParseString(v)
			if err != nil {
				return nil, err
			}
			resp = append(resp, value)
		}
		return resp, nil
	case []error:
		resp := make([]string, 0, len(val))
		for _, err := range val {
			resp = append(resp, err.Error())
		}
		return resp, nil
	default:
		return nil, ErrInvalidType
	}
}
