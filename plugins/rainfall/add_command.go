package rainfall

import (
	"fmt"
	"strconv"
)

type addCommand struct {
	p *plugin
}

func NewAddCommand(p *plugin) Commander {
	return &addCommand{
		p: p,
	}
}

func (c *addCommand) Execute(params []string) (string, error) {
	params = params[1:]
	loc, err := c.makeLocation(params)
	if err != nil {
		return "", err
	}

	_, ok := c.p.locStore.Get(loc.Name)
	if ok {
		return "", fmt.Errorf("%s は既に登録されています。", loc.Name)
	}

	c.p.locStore.Set(loc)

	return fmt.Sprintf("登録しました : %s [%f, %f]", loc.Name, loc.Latitude, loc.Longitude), nil
}

func (c *addCommand) makeLocation(params []string) (Location, error) {
	var loc Location

	if len(params) != 3 {
		return loc, CommandSyntaxError
	}

	loc.Name = params[0]

	lat, err := strconv.ParseFloat(params[1], 32)
	if err != nil {
		return loc, CommandSyntaxError
	}
	loc.Latitude = float32(lat)

	long, err := strconv.ParseFloat(params[2], 32)
	if err != nil {
		return loc, CommandSyntaxError
	}
	loc.Longitude = float32(long)

	return loc, nil
}
