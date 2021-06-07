package cast

import "fmt"

func Scan(tuple []interface{}, dest ...interface{}) error {
	if len(dest) > len(tuple) {
		return fmt.Errorf("got invalid destination size %d, tuple len=%d", len(tuple), len(dest))
	}

	var err error
	for i := 0; i < len(dest); i++ {
		switch t := dest[i].(type) {
		case *string:
			*t, err = ParseString(tuple[i])
		case *uint:
			*t, err = ParseUint(tuple[i])
		case *uint8:
			*t, err = ParseUint8(tuple[i])
		case *uint16:
			*t, err = ParseUint16(tuple[i])
		case *uint32:
			*t, err = ParseUint32(tuple[i])
		case *uint64:
			*t, err = ParseUint64(tuple[i])
		case *int:
			*t, err = ParseInt(tuple[i])
		case *int8:
			*t, err = ParseInt8(tuple[i])
		case *int16:
			*t, err = ParseInt16(tuple[i])
		case *int32:
			*t, err = ParseInt32(tuple[i])
		case *int64:
			*t, err = ParseInt64(tuple[i])
		case *float32:
			*t, err = ParseFloat32(tuple[i])
		case *float64:
			*t, err = ParseFloat64(tuple[i])
		default:
			err = fmt.Errorf("got unknown type %T", t)
		}
		if err != nil {
			err = fmt.Errorf("failed to parse %d destination value: %s", i, err)
			break
		}
	}
	return err
}
