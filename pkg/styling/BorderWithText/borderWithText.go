package BorderWithText

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// lipgloss doesn't really allow us to create a border like
//
// ╭── Some text here ──────╮
// │                        │
// ╰────────────────────────╯
//
// so we have to cheat slightly, we have to make a custom border object
// that will render the content of the inside of it

type Border struct {
	Top          string
	Bottom       string
	Left         string
	Right        string
	TopLeft      string
	TopRight     string
	BottomLeft   string
	BottomRight  string
	MiddleLeft   string
	MiddleRight  string
	Middle       string
	MiddleTop    string
	MiddleBottom string
}

func GetRoundTextBorder() Border {
	return Border{
		Top:          "─",
		Bottom:       "─",
		Left:         "│",
		Right:        "│",
		TopLeft:      "╭",
		TopRight:     "╮",
		BottomLeft:   "╰",
		BottomRight:  "╯",
		MiddleLeft:   "├",
		MiddleRight:  "┤",
		Middle:       "┼",
		MiddleTop:    "┬",
		MiddleBottom: "┴",
	}
}

// Apply the custom border to the content
// title will be what shows up as the border title
func ApplyBorder(border Border, title, content string, style Styles) string {

	lines, width := getLines(content)
	var out strings.Builder

	top := renderHorizontalEdge(border.TopLeft, border.Middle, border.TopRight, title, width, style)
	bottom := renderHorizontalEdge(border.BottomRight, border.Middle, border.BottomRight, "", width, style)

	out.WriteString(top)
	for i, l := range lines {
		out.WriteString(style.BorderStyle.Render(border.Left))
		out.WriteString(style.TextStyle.Render(l))
		out.WriteString(style.BorderStyle.Render(border.Right))
		if i < len(lines)-1 {
			out.WriteRune('\n')
		}
	}
	out.WriteString(bottom)

	return out.String()
}

func renderHorizontalEdge(left, middle, right, title string, width int, style Styles) string {
	leftWidth := ansi.StringWidth(left)
	rightWidth := ansi.StringWidth(left)
	middleWidth := ansi.StringWidth(middle)
	emptyWidth := ansi.StringWidth(" ")
	titleWidth := ansi.StringWidth(title)

	out := strings.Builder{}

	runes := []rune(middle)
	out.WriteString(style.BorderStyle.Render(left))

	spaceAlreadyTaken := leftWidth + rightWidth

	// if we've got the title, write it right away and update the length calculations
	if titleWidth > 0 {
		out.WriteString(style.BorderStyle.Render(middle))
		out.WriteString(style.BorderStyle.Render(middle))
		out.WriteString(style.TextStyle.Render(" "))
		out.WriteString(style.TextStyle.Render(title))
		out.WriteString(style.TextStyle.Render(" "))
		spaceAlreadyTaken = spaceAlreadyTaken + titleWidth + (2 * middleWidth) + (2 * emptyWidth)
	}

	j := 0
	runesLen := len(runes)
	for i := spaceAlreadyTaken; i < width+rightWidth; {
		out.WriteString(style.BorderStyle.Render(string(runes[j])))
		j++
		if j >= runesLen {
			j = 0
		}

		i += ansi.StringWidth((string(runes[j])))
	}

	out.WriteString(style.BorderStyle.Render(right))
	return out.String()
}

func getLines(s string) (lines []string, widest int) {
	lines = strings.Split(s, "\n")

	for _, l := range lines {
		w := ansi.StringWidth(l)
		if widest < w {
			widest = w
		}
	}

	return lines, widest
}
