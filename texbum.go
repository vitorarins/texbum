package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

// TcSetAttr restores the terminal connected to the given file descriptor to a
// previous state.
func TcSetAttr(fd uintptr, termios *syscall.Termios) error {
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(termios))); err != 0 {
		return err
	}
	return nil
}

// TcGetAttr retrieves the current terminal settings and returns it.
func TcGetAttr(fd uintptr, termios *syscall.Termios) error {
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios))); err != 0 {
		return err
	}
	return nil
}

func enableRawMode() {
	out, err := os.OpenFile("/dev/tty", syscall.O_WRONLY, 0)
	if err != nil {
		fmt.Printf("Could not open TTY because: %s", err)
	}
	
	var raw syscall.Termios
	err = TcGetAttr(out.Fd(), &raw)
	
	if err != nil {
		fmt.Printf("Could not GET attribute from terminal because: %s\n", err)
	}

	raw.Lflag &^= syscall.ECHO;

	err = TcSetAttr(out.Fd(), &raw)

	if err != nil {
		fmt.Println("Could not SET attribute from terminal.")
	}
}

func main() {
	enableRawMode()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan()  {
		name := scanner.Text()
		if strings.Index(name,"q") > -1 {
			break;
		}
		fmt.Println(name)
	}
}
