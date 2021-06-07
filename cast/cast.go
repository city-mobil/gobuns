package cast

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var (
	ErrNegativeToUnsigned = errors.New("unable to cast from negative to unsigned value")
	ErrInvalidType        = errors.New("unable to cast value: invalid type")
)

func ParseInt(f interface{}) (int, error) {
	switch t := f.(type) {
	case uint:
		return int(t), nil
	case int:
		return t, nil
	case uint8:
		return int(t), nil
	case int8:
		return int(t), nil
	case uint16:
		return int(t), nil
	case int16:
		return int(t), nil
	case uint32:
		return int(t), nil
	case int32:
		return int(t), nil
	case uint64:
		return int(t), nil
	case int64:
		return int(t), nil
	case string:
		return strconv.Atoi(t)
	}
	return 0, ErrInvalidType
}

func ParseUint(i interface{}) (uint, error) {
	switch v := i.(type) {
	case uint:
		return v, nil
	case int:
		if v < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint(v), nil
	case uint8:
		return uint(v), nil
	case int8:
		if v < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint(v), nil
	case uint16:
		return uint(v), nil
	case int16:
		if v < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint(v), nil
	case uint32:
		return uint(v), nil
	case int32:
		if v < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint(v), nil
	case uint64:
		return uint(v), nil
	case int64:
		if v < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint(v), nil
	case string:
		val, err := strconv.ParseUint(v, 10, 0)
		return uint(val), err
	}
	return 0, ErrInvalidType
}

func ParseUint8(f interface{}) (uint8, error) {
	switch t := f.(type) {
	case uint:
		return uint8(t), nil
	case int:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint8(t), nil
	case uint8:
		return t, nil
	case int8:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint8(t), nil
	case uint16:
		return uint8(t), nil
	case int16:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint8(t), nil
	case uint32:
		return uint8(t), nil
	case int32:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint8(t), nil
	case uint64:
		return uint8(t), nil
	case int64:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint8(t), nil
	case string:
		v, err := strconv.ParseUint(t, 10, 8)
		return uint8(v), err
	}
	return 0, ErrInvalidType
}

func ParseUint16(f interface{}) (uint16, error) {
	switch t := f.(type) {
	case uint:
		return uint16(t), nil
	case int:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint16(t), nil
	case uint8:
		return uint16(t), nil
	case int8:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint16(t), nil
	case uint16:
		return t, nil
	case int16:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint16(t), nil
	case uint32:
		return uint16(t), nil
	case int32:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint16(t), nil
	case uint64:
		return uint16(t), nil
	case int64:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint16(t), nil
	case string:
		v, err := strconv.ParseUint(t, 10, 16)
		return uint16(v), err
	}
	return 0, ErrInvalidType
}

func ParseUint32(f interface{}) (uint32, error) {
	switch t := f.(type) {
	case uint:
		return uint32(t), nil
	case int:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint32(t), nil
	case uint8:
		return uint32(t), nil
	case int8:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint32(t), nil
	case uint16:
		return uint32(t), nil
	case int16:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint32(t), nil
	case uint32:
		return t, nil
	case int32:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint32(t), nil
	case uint64:
		return uint32(t), nil
	case int64:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint32(t), nil
	case string:
		v, err := strconv.ParseUint(t, 10, 32)
		return uint32(v), err
	}
	return 0, ErrInvalidType
}

func ParseUint64(f interface{}) (uint64, error) {
	switch t := f.(type) {
	case uint:
		return uint64(t), nil
	case int:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint64(t), nil
	case uint8:
		return uint64(t), nil
	case int8:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint64(t), nil
	case int16:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint64(t), nil
	case uint16:
		return uint64(t), nil
	case int32:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint64(t), nil
	case uint32:
		return uint64(t), nil
	case uint64:
		return t, nil
	case int64:
		if t < 0 {
			return 0, ErrNegativeToUnsigned
		}
		return uint64(t), nil
	case string:
		return strconv.ParseUint(t, 10, 64)
	default:
		return 0, ErrInvalidType
	}
}

func ParseInt8(f interface{}) (int8, error) {
	switch t := f.(type) {
	case uint:
		return int8(t), nil
	case int:
		return int8(t), nil
	case uint8:
		return int8(t), nil
	case int8:
		return t, nil
	case uint16:
		return int8(t), nil
	case int16:
		return int8(t), nil
	case uint32:
		return int8(t), nil
	case int32:
		return int8(t), nil
	case uint64:
		return int8(t), nil
	case int64:
		return int8(t), nil
	case string:
		v, err := strconv.ParseInt(t, 10, 8)
		return int8(v), err
	}
	return 0, ErrInvalidType
}

func ParseInt16(f interface{}) (int16, error) {
	switch t := f.(type) {
	case uint:
		return int16(t), nil
	case int:
		return int16(t), nil
	case uint8:
		return int16(t), nil
	case int8:
		return int16(t), nil
	case uint16:
		return int16(t), nil
	case int16:
		return t, nil
	case uint32:
		return int16(t), nil
	case int32:
		return int16(t), nil
	case uint64:
		return int16(t), nil
	case int64:
		return int16(t), nil
	case string:
		v, err := strconv.ParseInt(t, 10, 16)
		return int16(v), err
	}
	return 0, ErrInvalidType
}

func ParseInt32(f interface{}) (int32, error) {
	switch t := f.(type) {
	case uint:
		return int32(t), nil
	case int:
		return int32(t), nil
	case uint8:
		return int32(t), nil
	case int8:
		return int32(t), nil
	case uint16:
		return int32(t), nil
	case int16:
		return int32(t), nil
	case uint32:
		return int32(t), nil
	case int32:
		return t, nil
	case uint64:
		return int32(t), nil
	case int64:
		return int32(t), nil
	case string:
		v, err := strconv.ParseInt(t, 10, 32)
		return int32(v), err
	}
	return 0, ErrInvalidType
}

func ParseInt64(f interface{}) (int64, error) {
	switch t := f.(type) {
	case int:
		return int64(t), nil
	case uint:
		return int64(t), nil
	case uint8:
		return int64(t), nil
	case int8:
		return int64(t), nil
	case uint16:
		return int64(t), nil
	case int16:
		return int64(t), nil
	case uint32:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case uint64:
		return int64(t), nil
	case int64:
		return t, nil
	case string:
		return strconv.ParseInt(t, 10, 64)
	}
	return 0, ErrInvalidType
}

func ParseFloat32(f interface{}) (float32, error) {
	switch t := f.(type) {
	case uint:
		return float32(t), nil
	case int:
		return float32(t), nil
	case uint8:
		return float32(t), nil
	case int8:
		return float32(t), nil
	case uint16:
		return float32(t), nil
	case int16:
		return float32(t), nil
	case uint32:
		return float32(t), nil
	case int32:
		return float32(t), nil
	case uint64:
		return float32(t), nil
	case int64:
		return float32(t), nil
	case float32:
		return t, nil
	case float64:
		return float32(t), nil
	case string:
		v, err := strconv.ParseFloat(t, 32)
		return float32(v), err
	}

	return 0, ErrInvalidType
}

func ParseFloat64(f interface{}) (float64, error) {
	switch t := f.(type) {
	case uint:
		return float64(t), nil
	case int:
		return float64(t), nil
	case uint8:
		return float64(t), nil
	case int8:
		return float64(t), nil
	case uint16:
		return float64(t), nil
	case int16:
		return float64(t), nil
	case uint32:
		return float64(t), nil
	case int32:
		return float64(t), nil
	case uint64:
		return float64(t), nil
	case int64:
		return float64(t), nil
	case float32:
		return float64(t), nil
	case float64:
		return t, nil
	case string:
		v, err := strconv.ParseFloat(t, 64)
		return v, err
	}

	return 0, ErrInvalidType
}

func ParseString(f interface{}) (string, error) {
	switch t := f.(type) {
	case string:
		return t, nil
	case bool:
		return strconv.FormatBool(t), nil
	case float32:
		return strconv.FormatFloat(float64(t), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(t), nil
	case uint:
		return strconv.FormatUint(uint64(t), 10), nil
	case int8:
		return strconv.FormatInt(int64(t), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(t), 10), nil
	case int16:
		return strconv.FormatInt(int64(t), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(t), 10), nil
	case int32:
		return strconv.Itoa(int(t)), nil
	case uint32:
		return strconv.FormatUint(uint64(t), 10), nil
	case int64:
		return strconv.FormatInt(t, 10), nil
	case uint64:
		return strconv.FormatUint(t, 10), nil
	case []byte:
		return string(t), nil
	case nil:
		return "", nil
	case fmt.Stringer:
		return t.String(), nil
	case error:
		return t.Error(), nil
	}

	return "", ErrInvalidType
}

func ParseBool(f interface{}) (bool, error) {
	switch t := f.(type) {
	case bool:
		return t, nil
	case nil:
		return false, nil
	case string:
		return strconv.ParseBool(t)
	case uint:
		return t != uint(0), nil
	case int:
		return t != 0, nil
	case uint8:
		return t != uint8(0), nil
	case int8:
		return t != int8(0), nil
	case uint16:
		return t != uint16(0), nil
	case int16:
		return t != int16(0), nil
	case uint32:
		return t != uint32(0), nil
	case int32:
		return t != int32(0), nil
	case uint64:
		return t != uint64(0), nil
	case int64:
		return t != int64(0), nil
	case float32:
		return t != float32(0), nil
	case float64:
		return t != float64(0), nil
	}

	return false, ErrInvalidType
}

func ParseDuration(f interface{}) (time.Duration, error) {
	switch t := f.(type) {
	case time.Duration:
		return t, nil
	case float32, float64:
		v, err := ParseFloat64(t)
		if err != nil {
			return 0, err
		}
		return time.Duration(v), nil
	case string:
		v, err := time.ParseDuration(t)
		if err != nil {
			return 0, err
		}
		return v, err
	default:
		v, err := ParseInt64(t)
		if err != nil {
			return 0, err
		}
		return time.Duration(v), nil
	}
}
