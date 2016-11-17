package rainfall

import (
	"bytes"
	"fmt"
)

type listCommand struct {
	p *plugin
}

func NewListCommand(p *plugin) Commander {
	return &listCommand{
		p: p,
	}
}

func (c *listCommand) Execute(params []string) (string, error) {
	params = params[1:]
	if len(params) != 0 {
		return "", CommandSyntaxError
	}

	locations := c.p.locStore.Locations()

	resMessage := bytes.Buffer{}
	for i, loc := range locations {
		if i > 0 {
			resMessage.WriteString("\n")
		}
		resMessage.WriteString(fmt.Sprintf("%s [%f, %f]", loc.Name, loc.Latitude, loc.Longitude))
	}

	return resMessage.String(), nil
}
