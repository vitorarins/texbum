package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var origTermios syscall.Termios

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

func disableRawMode() {
	out, err := os.OpenFile("/dev/tty", syscall.O_WRONLY, 0)
	defer out.Close()

	err = TcSetAttr(out.Fd(), &origTermios)

	if err != nil {
		fmt.Printf("Could not SET attribute from original terminal because: %s.\n", err)
	}
}

func enableRawMode() {
	out, err := os.OpenFile("/dev/tty", syscall.O_WRONLY, 0)
	defer out.Close()
	if err != nil {
		fmt.Printf("Could not open TTY because: %s", err)
	}
	
	err = TcGetAttr(out.Fd(), &origTermios)
	
	if err != nil {
		fmt.Printf("Could not GET attribute from terminal because: %s.\n", err)
	}

	raw := origTermios
	raw.Lflag &^= syscall.ECHO;

	err = TcSetAttr(out.Fd(), &raw)

	if err != nil {
		fmt.Printf("Could not SET attribute from terminal because: %s.\n", err)
	}
}

func main() {
	enableRawMode()
	defer disableRawMode()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan()  {
		name := scanner.Text()
		if strings.Index(name,"q") > -1 {
			break;
		}
		fmt.Println(name)
	}
}
