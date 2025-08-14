package simplelog

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// msg represets fields of log message
type msg struct {
	TimeStamp string
	Prefix    string
	Text      string
}

// String return string representation of message
func (m *msg) String() string {
	sb := new(strings.Builder)

	if m.TimeStamp != "" {
		sb.WriteString(m.TimeStamp)
		sb.WriteRune(' ')
	}

	if m.Prefix != "" {
		sb.WriteString(m.Prefix)
		sb.WriteRune(' ')
	}

	sb.WriteString(m.Text)

	return sb.String()
}

// fit fits whole message to specified width `width` by reducing message text if needed. If message does not fit needed
// width trim marker `trimMarker` will added to the end of message text.
func (m *msg) fit(width int, trimMarker string) {
	spaceCount := 0
	if m.TimeStamp != "" {
		spaceCount = 1
	}

	spaceLeft := width - lipgloss.Width(m.TimeStamp) - lipgloss.Width(m.Prefix) - lipgloss.Width(m.Text) - spaceCount
	if spaceLeft >= 0 {
		return
	}

	maxMessageWidth := max(len([]rune(m.Text))+spaceLeft-len(trimMarker), 0)

	m.Text = string([]rune(m.Text)[:maxMessageWidth]) + trimMarker
}
