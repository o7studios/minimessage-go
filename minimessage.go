package minimessage_go

import (
	"fmt"
	"go.minekube.com/common/minecraft/color"
	c "go.minekube.com/common/minecraft/component"
	key2 "go.minekube.com/common/minecraft/key"
	"math"
	"strings"
)

func Deserialize(input string) (*c.Text, error) {
	styleStack := []c.Style{{Color: color.White}}

	var components []c.Component

	if !strings.Contains(input, "<") || !strings.Contains(input, ">") {
		return &c.Text{
			Content: input,
			S:       styleStack[0],
		}, nil
	}

	for _, str := range strings.Split(input, "<") {
		if str == "" {
			continue
		}

		splitStr := strings.Split(str, ">")

		if len(splitStr) < 2 {
			continue
		}

		key := splitStr[0]

		if strings.HasPrefix(key, "/") {
			styleStack = styleStack[:len(styleStack)-1]
			key = key[1:]
		} else {
			newStyle := styleStack[len(styleStack)-1]

			styleStack = append(styleStack, newStyle)

			newText := modify(key, splitStr[1], &styleStack[len(styleStack)-1])

			if newText == nil {
				continue
			}

			components = append(components, newText)
		}
	}

	return &c.Text{
		Extra: components,
	}, nil
}

func modify(key string, content string, s *c.Style) *c.Text {
	newText := &c.Text{}

	switch {
	case strings.HasPrefix(key, "font"):
		fontKey := strings.Split(key, ":")
		fontNames := fontKey[1:]

		if len(fontNames) > 2 || len(fontNames) < 1 {
			return nil
		}

		if len(fontNames) < 2 {
			parsedFont, err := key2.Parse(fontNames[0])
			if err != nil {
				return nil
			}
			s.Font = parsedFont

			return newText
		}

		join := fontNames[0] + ":" + fontNames[1]
		parsedFont, err := key2.Parse(join)
		if err != nil {
			return nil
		}

		s.Font = parsedFont

		return newText
	case strings.HasPrefix(key, "gradient"):
		colorKey := strings.Split(key, ":")
		colorNames := colorKey[1:]

		colors := make([]color.RGB, len(colorNames))
		for i, col := range colorNames {
			parsedColor, err := parseColor(col)
			if err != nil {
				return nil
			}
			newColor, _ := color.Make(parsedColor)
			colors[i] = *newColor
		}

		newText = gradient(content, *s, colors...)
		return newText
	case strings.HasPrefix(key, "color:") || strings.HasPrefix(key, "colour:") || strings.HasPrefix(key, "c:"):
		key = strings.Split(key, ":")[1]
		fallthrough
	case strings.EqualFold(key, "bold") || strings.EqualFold(key, "b"):
		s.Bold = c.True
	case strings.EqualFold(key, "obfuscated") || strings.EqualFold(key, "obf"):
		s.Obfuscated = c.True
	case strings.EqualFold(key, "italic") || strings.EqualFold(key, "i") || strings.EqualFold(key, "em"):
		s.Italic = c.True
	case strings.EqualFold(key, "strikethrough") || strings.EqualFold(key, "st"):
		s.Strikethrough = c.True
	case strings.EqualFold(key, "underlined") || strings.EqualFold(key, "underline") || strings.EqualFold(key, "u"):
		s.Underlined = c.True
	case strings.EqualFold(key, "newline") || strings.EqualFold(key, "br"):
		content = "\n" + content
	default:
		parsed, err := parseColor(key)
		if err != nil {
			return nil
		}
		s.Color = parsed
	}

	newText.Content = content
	newText.S = *s

	return newText
}

func parseColor(name string) (color.Color, error) {
	if strings.HasPrefix(name, "#") {
		newColor, err := color.Hex(name)
		if err != nil {
			return nil, err
		}
		return newColor, nil
	}
	return parseColorFromName(name)
}

func parseColorFromName(name string) (color.Color, error) {
	col, ok := color.Names[name]
	if ok {
		return col, nil
	}
	for _, someName := range color.Names {
		if !strings.EqualFold(someName.String(), name) {
			continue
		}
		return someName, nil
	}
	return nil, fmt.Errorf("unknown color name: %s", name)
}

func gradient(content string, style c.Style, colors ...color.RGB) *c.Text {
	var component []c.Component
	for i, char := range strings.Split(content, "") {
		t := float64(i) / float64(len(content))
		hex, _ := color.Hex(lerpColor(t, colors...).Hex())

		style.Color = hex
		component = append(component, &c.Text{
			Content: char,
			S:       style,
		})
	}

	return &c.Text{
		Extra: component,
	}
}

func lerpColor(t float64, colors ...color.RGB) color.Color {
	t = math.Min(t, 1)

	if t == 1 {
		return &colors[len(colors)-1]
	}

	colorT := t * float64(len(colors)-1)
	newT := colorT - math.Floor(colorT)
	lastColor := colors[int(colorT)]
	nextColor := colors[int(colorT+1)]

	return &color.RGB{
		R: lerpInt(newT, nextColor.R, lastColor.R),
		G: lerpInt(newT, nextColor.G, lastColor.G),
		B: lerpInt(newT, nextColor.B, lastColor.B),
	}
}

func lerpInt(t float64, a float64, b float64) float64 {
	return a*t + b*(1-t)
}
