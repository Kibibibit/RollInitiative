package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

const ANSI_ESCAPE = "\x1b"
const ANSI_SEQ = ANSI_ESCAPE + "["

const ANSI_ERASE_SCREEN = ANSI_SEQ + "2J"

const ANSI_HIDE_CURSOR = ANSI_SEQ + "?25l"
const ANSI_SHOW_CURSOR = ANSI_SEQ + "?25h"

const ANSI_SAVE_SCREEN = ANSI_SEQ + "?47h"
const ANSI_RESTORE_SCREEN = ANSI_SEQ + "?47l"

const ioctlReadTermios = unix.TCGETS
const ioctlWriteTermios = unix.TCSETS

type TTY struct {
	in         *os.File
	bin        *bufio.Reader
	out        *os.File
	termios    unix.Termios
	ss         chan os.Signal
	size       IVector2
	drawBuffer string
}

func Open() (*TTY, error) {
	return open("/dev/tty")
}

func open(path string) (*TTY, error) {
	tty := new(TTY)

	in, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	tty.in = in
	tty.bin = bufio.NewReader(in)

	out, err := os.OpenFile(path, syscall.O_WRONLY, 0)
	if err != nil {
		return nil, err
	}
	tty.out = out

	termios, err := unix.IoctlGetTermios(int(tty.in.Fd()), ioctlReadTermios)
	if err != nil {
		return nil, err
	}

	tty.termios = *termios

	termios.Iflag &^= unix.ISTRIP | unix.INLCR | unix.ICRNL | unix.IGNCR | unix.IXOFF
	termios.Lflag &^= unix.ECHO | unix.ICANON /*| unix.ISIG*/
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0
	if err := unix.IoctlSetTermios(int(tty.in.Fd()), ioctlWriteTermios, termios); err != nil {
		return nil, err
	}

	tty.ss = make(chan os.Signal, 1)

	signal.Notify(tty.ss, os.Interrupt)

	x, y, err := term.GetSize(int(tty.in.Fd()))

	if err != nil {
		return nil, err
	}

	tty.size = IVector2{x: x, y: y}

	tty.out.WriteString(ANSI_SAVE_SCREEN)
	tty.ClearScreen()
	tty.HideCursor()

	tty.drawBuffer = ""

	return tty, nil
}

func (tty *TTY) Close() error {
	tty.out.WriteString(ANSI_RESTORE_SCREEN)
	tty.ShowCursor()
	if tty.out == nil || tty.in == nil {
		return nil
	}
	signal.Stop(tty.ss)
	close(tty.ss)
	ioctlErr := unix.IoctlSetTermios(int(tty.in.Fd()), ioctlWriteTermios, &tty.termios)
	outErr := tty.out.Close()
	inErr := tty.in.Close()

	tty.out = nil
	tty.in = nil

	if ioctlErr != nil {
		log.Fatal("IOCTL ERROR")
		return ioctlErr
	}
	if outErr != nil {
		return outErr
	}
	return inErr

}

func (tty *TTY) Reset() string {
	return fmt.Sprintf("%s0m", ANSI_SEQ)
}

func (tty *TTY) ReadRune() (rune, int, error) {
	r, size, err := tty.bin.ReadRune()
	return r, size, err
}

func (tty *TTY) DrawString(data string) {
	tty.out.WriteString(data)
}

func (tty *TTY) ClearScreen() {
	tty.out.WriteString(ANSI_ERASE_SCREEN)
}

func (tty *TTY) HideCursor() {
	tty.out.WriteString(ANSI_HIDE_CURSOR)
}

func (tty *TTY) ShowCursor() {
	tty.out.WriteString(ANSI_SHOW_CURSOR)
}

func (tty *TTY) MoveCursorString(x int, y int) string {
	return fmt.Sprintf("%s%d;%dH", ANSI_SEQ, y, x)
}

func (tty *TTY) ForegroundRGBColorString(r int, g int, b int) string {
	return fmt.Sprintf("%s38;2;%d;%d;%dm", ANSI_SEQ, r, g, b)
}

func (tty *TTY) BackgroundRGBColorString(r int, g int, b int) string {
	return fmt.Sprintf("%s48;2;%d;%d;%dm", ANSI_SEQ, r, g, b)
}

func (tty *TTY) AddToBuffer(data string) {
	tty.drawBuffer += data
}

func (tty *TTY) FlushBuffer() {
	tty.out.WriteString(tty.drawBuffer)
	tty.drawBuffer = ""
}
