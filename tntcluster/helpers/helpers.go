package helpers

import (
	"github.com/viciious/go-tarantool"
)

type TarantoolCommand string

const (
	SelectCommand  TarantoolCommand = "Tarantool: Select"
	CallCommand    TarantoolCommand = "Tarantool: Call"
	AuthCommand    TarantoolCommand = "Tarantool: Auth"
	InsertCommand  TarantoolCommand = "Tarantool: Insert"
	ReplaceCommand TarantoolCommand = "Tarantool: Replace"
	DeleteCommand  TarantoolCommand = "Tarantool: Delete"
	UpdateCommand  TarantoolCommand = "Tarantool: Update"
	UpsertCommand  TarantoolCommand = "Tarantool: Upsert"
	PingCommand    TarantoolCommand = "Tarantool: Ping"
	EvalCommand    TarantoolCommand = "Tarantool: Eval"
	UnknownCommand TarantoolCommand = "Tarantool: Unknown"
)

var cmdName = map[uint]TarantoolCommand{
	tarantool.SelectCommand:  SelectCommand,
	tarantool.AuthCommand:    AuthCommand,
	tarantool.InsertCommand:  InsertCommand,
	tarantool.ReplaceCommand: ReplaceCommand,
	tarantool.DeleteCommand:  DeleteCommand,
	tarantool.UpdateCommand:  UpdateCommand,
	tarantool.UpsertCommand:  UpsertCommand,
	tarantool.PingCommand:    PingCommand,
	tarantool.EvalCommand:    EvalCommand,
}

// TarantoolCommandAndStatement returns tarantool operation name and statement corresponding to the command.
func TarantoolCommandAndStatement(query tarantool.Query) (cmd TarantoolCommand, statement string) {
	cmd, ok := cmdName[query.GetCommandID()]
	if !ok {
		cmd = UnknownCommand
	}

	if query.GetCommandID() == tarantool.CallCommand {
		cmd = CallCommand
		if call, ok := query.(*tarantool.Call); ok {
			statement = call.Name
		}
	}
	return
}
