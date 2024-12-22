package minimessage_go

import (
	"go.minekube.com/common/minecraft/color"
	c "go.minekube.com/common/minecraft/component"
	"strings"
)

func deserialize(input string) (*c.Text, error) {
	styleStack := []c.Style{{Color: color.White}}

	var components []c.Component

	if !strings.Contains(input, "<") || !strings.Contains(input, ">") {
		return &c.Text{
			Content: input,
			S:       styleStack[0],
		}, nil
	}

	return &c.Text{
		Extra: components,
	}, nil
}
