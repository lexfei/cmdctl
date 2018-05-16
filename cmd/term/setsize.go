// +build !windows

package term

import (
	"cmdctl/pkg/term"
)

// SetSize sets the terminal size associated with fd.
func SetSize(fd uintptr, size TerminalSize) error {
	return term.SetWinsize(fd, &term.Winsize{Height: size.Height, Width: size.Width})
}
