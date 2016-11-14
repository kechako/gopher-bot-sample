package rainfall

import (
	"errors"
	"strings"
)

type CommandType string

const (
	CommandTypeAddLocation    CommandType = "add"
	CommandTypeRemoveLocation CommandType = "rm"
	CommandTypeChangeLocation CommandType = "change"
	CommandTypeNone           CommandType = "none"
)

var (
	CommandSyntaxError = errors.New("CommandSyntaxError")
)

type Commander interface {
	Execute(params []string) (string, error)
}

func (p *plugin) ExecuteCommand(message string) (string, error) {
	params := make([]string, 0, 4)
	for _, param := range strings.Split(message, " ") {
		if param == "" {
			continue
		}
		params = append(params, param)
	}

	cmdType := CommandTypeNone
	if len(params) > 0 {
		cmdType = CommandType(params[0])
	}

	var commander Commander
	switch cmdType {
	case CommandTypeAddLocation:
	case CommandTypeRemoveLocation:
	case CommandTypeChangeLocation:
	default:
		commander = NewAskCommand(p)
	}

	return commander.Execute(params)
}