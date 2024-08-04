package gotestdatabase

import (
	"fmt"
	"strings"
)

type TreeStyler interface {
	getLeft() string
	getRight() string
	getLine() string
	getFork() string
	getLast() string
}

type TreeFrame struct {
	lines  []string
	styler TreeStyler
}

func (tf *TreeFrame) render(root *DatabaseNode, render bool) string {
	root.render(tf, "", "", -1)

	lines := make([]string, 0, len(tf.lines))
	for k, v := range tf.lines {
		runeV := []rune(v)
		if string(runeV[0]) != " " {
			lines = append(lines, v)
			continue
		}

		line := ""
		for i, j := range runeV {
			previousLine := []rune(tf.lines[k-1])
			if i >= len(previousLine) {
				line += string(j)
				continue
			}

			previousChar := fmt.Sprintf("%c", previousLine[i])
			previousNewChar := fmt.Sprintf("%c", []rune(lines[k-1])[i])
			if (previousChar == tf.styler.getFork() || previousNewChar == tf.styler.getLine()) && fmt.Sprintf("%c", j) == " " {
				line += tf.styler.getLine()
			} else {
				line += string(j)
			}
		}

		lines = append(lines, line)
	}

	renderable := strings.Join(lines, "\n")
	if render {
		print(renderable)
	}

	return renderable
}
