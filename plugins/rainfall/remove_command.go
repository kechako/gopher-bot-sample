package rainfall

import "fmt"

type removeCommand struct {
	p *plugin
}

func NewRemoveCommand(p *plugin) Commander {
	return &removeCommand{
		p: p,
	}
}

func (c *removeCommand) Execute(params []string) (string, error) {
	params = params[1:]
	if len(params) != 1 {
		return "", CommandSyntaxError
	}

	name := params[0]

	loc, ok := c.p.locStore.Get(name)
	if !ok {
		return "", fmt.Errorf("%s は登録されていません。", name)
	}

	c.p.locStore.Del(name)

	return fmt.Sprintf("削除しました : %s [%f, %f]", loc.Name, loc.Latitude, loc.Longitude), nil
}
