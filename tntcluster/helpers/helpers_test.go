package helpers

import (
	"testing"

	"github.com/viciious/go-tarantool"
)

func TestTarantoolCommandAndStatement(t *testing.T) {
	var testData = []struct {
		command           tarantool.Query
		expectedStatement string
		expectedCmd       TarantoolCommand
		testName          string
	}{
		{
			command: &tarantool.Call{
				Name:  "testName",
				Tuple: []interface{}{},
			},
			expectedStatement: "testName",
			expectedCmd:       CallCommand,
			testName:          "call",
		},
		{
			command:     &tarantool.Select{},
			expectedCmd: SelectCommand,
			testName:    "select",
		},
	}

	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			cmd, st := TarantoolCommandAndStatement(v.command)
			if st != v.expectedStatement {
				t.Errorf("got statement %q, expected %q", st, v.expectedStatement)
			}
			if cmd != v.expectedCmd {
				t.Errorf("got cmd %q, expected %q", cmd, v.expectedCmd)
			}
		})
	}
}
