package rainfall

import (
	"fmt"
	"strconv"
)

type changeCommand struct {
	p *plugin
}

func NewChangeCommand(p *plugin) Commander {
	return &changeCommand{
		p: p,
	}
}

func (c *changeCommand) Execute(params []string) (string, error) {
	params = params[1:]
	loc, err := c.makeLocation(params)
	if err != nil {
		return "", err
	}

	oldLoc, ok := c.p.locStore.Get(loc.Name)
	if !ok {
		return "", fmt.Errorf("%s は登録されていません。", loc.Name)
	}

	c.p.locStore.Set(loc)

	return fmt.Sprintf("変更しました : %s [%f, %f] => %s [%f, %f]", oldLoc.Name, oldLoc.Latitude, oldLoc.Longitude, loc.Name, loc.Latitude, loc.Longitude), nil
}

func (c *changeCommand) makeLocation(params []string) (Location, error) {
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
