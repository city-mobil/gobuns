package mysql

import (
	"database/sql/driver"
	"net"

	"github.com/go-sql-driver/mysql"
)

var (
	retryableErrors = map[uint16]struct{}{
		1040: {}, // ER_CON_COUNT_ERROR
		1042: {}, // ER_BAD_HOST_ERROR
		1043: {}, // ER_HANDSHAKE_ERROR
		1053: {}, // ER_SERVER_SHUTDOWN
		1317: {}, // ER_QUERY_INTERRUPTED
	}
)

// ErrorCode returns the MySQL server error code for the error, or zero
// if the error is not a MySQL error.
func ErrorCode(err error) uint16 {
	if val, ok := err.(*mysql.MySQLError); ok {
		return val.Number
	}
	return 0 // not a mysql error
}

// CanRetry returns true for every error which can be safely retry.
// It returns false for all other errors, including nil.
func CanRetry(err error) bool {
	if err == nil {
		return false
	}

	if err == mysql.ErrInvalidConn || err == driver.ErrBadConn {
		return true
	}

	if _, ok := err.(*net.OpError); ok {
		// Being unable to reach MySQL is a network issue, so we get a net.OpError.
		// If MySQL is reachable, then we'd get a mysql.* or driver.* error instead.
		return true
	}

	code := ErrorCode(err)
	_, ok := retryableErrors[code]
	return ok
}
