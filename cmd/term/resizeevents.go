// +build !windows

package term

import (
	"os"
	"os/signal"
	"syscall"

	"cmdctl/pkg/runtime"
)

// monitorResizeEvents spawns a goroutine that waits for SIGWINCH signals (these indicate the
// terminal has resized). After receiving a SIGWINCH, this gets the terminal size and tries to send
// it to the resizeEvents channel. The goroutine stops when the stop channel is closed.
func monitorResizeEvents(fd uintptr, resizeEvents chan<- TerminalSize, stop chan struct{}) {
	go func() {
		defer runtime.HandleCrash()

		winch := make(chan os.Signal, 1)
		signal.Notify(winch, syscall.SIGWINCH)
		defer signal.Stop(winch)

		for {
			select {
			case <-winch:
				size := GetSize(fd)
				if size == nil {
					return
				}

				// try to send size
				select {
				case resizeEvents <- *size:
					// success
				default:
					// not sent
				}
			case <-stop:
				return
			}
		}
	}()
}
