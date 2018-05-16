// +build windows

package term

func SetSize(fd uintptr, size TerminalSize) error {
	// NOP
	return nil
}
