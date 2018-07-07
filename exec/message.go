package exec

import (
	"fmt"

	"github.com/mickep76/runshit/color"
)

func (m *Message) FmtString() string {
	return fmt.Sprintf("%s\n", m.Message)
}

func (m *Message) FmtStringColor(col int, addr string) string {
	return fmt.Sprintf("%s%s%s %s\n", color.Fg256(uint8(col+1)), addr, color.Reset, m.Message)
}
