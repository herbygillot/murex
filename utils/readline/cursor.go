package readline

import (
	"os"
	"regexp"
	"strconv"
)

func leftMost() []byte {
	fd := int(os.Stdout.Fd())
	w, _, err := GetSize(fd)
	if err != nil {
		return []byte{'\r', '\n'}
	}

	b := make([]byte, w+1)
	for i := 0; i < w; i++ {
		b[i] = ' '
	}
	b[w] = '\r'

	return b
}

var rxRcvCursorPos = regexp.MustCompile("^\x1b([0-9]+);([0-9]+)R$")

func (rl *Instance) getCursorPos() (x int, y int) {
	if !rl.EnableGetCursorPos {
		return -1, -1
	}

	disable := func() (int, int) {
		printErr("\r\ngetCursorPos() not supported by terminal emulator, disabling....\r\n")
		rl.EnableGetCursorPos = false
		return -1, -1
	}

	print(seqGetCursorPos)
	b := make([]byte, 64)
	i, err := os.Stdin.Read(b)
	if err != nil {
		return disable()
	}

	if !rxRcvCursorPos.Match(b[:i]) {
		return disable()
	}

	match := rxRcvCursorPos.FindAllStringSubmatch(string(b[:i]), 1)
	y, err = strconv.Atoi(match[0][1])
	if err != nil {
		return disable()
	}

	x, err = strconv.Atoi(match[0][2])
	if err != nil {
		return disable()
	}

	return x, y
}

func moveCursorUp(i int) {
	if i < 1 {
		return
	}

	printf("\x1b[%dA", i)
}

func moveCursorDown(i int) {
	if i < 1 {
		return
	}

	printf("\x1b[%dB", i)
}

func moveCursorForwards(i int) {
	if i < 1 {
		return
	}

	printf("\x1b[%dC", i)
}

func moveCursorBackwards(i int) {
	if i < 1 {
		return
	}

	printf("\x1b[%dD", i)
}

func (rl *Instance) moveCursorToStart() {
	posX, posY := lineWrapPos(rl.promptLen, rl.pos, rl.termWidth)

	moveCursorBackwards(posX - rl.promptLen)
	moveCursorUp(posY)
}

func (rl *Instance) moveCursorFromStartToLinePos() {
	posX, posY := lineWrapPos(rl.promptLen, rl.pos, rl.termWidth)
	moveCursorForwards(posX)
	moveCursorDown(posY)
}

func (rl *Instance) moveCursorFromEndToLinePos() {
	lineX, lineY := lineWrapPos(rl.promptLen, len(rl.line), rl.termWidth)
	posX, posY := lineWrapPos(rl.promptLen, rl.pos, rl.termWidth)
	moveCursorBackwards(lineX - posX)
	moveCursorUp(lineY - posY)
}

/*// moveCursorToLinePos should only be used on extreme circumstances because it
// causes the cursor to jump around quite a bit
func (rl *Instance) moveCursorFromUnknownToLinePos() {
	_, lineY := lineWrapPos(rl.promptLen, len(rl.line), rl.termWidth)
	posX, posY := lineWrapPos(rl.promptLen, rl.pos, rl.termWidth)
	//moveCursorBackwards(lineX)
	print("\r")
	moveCursorForwards(posX)
	moveCursorUp(lineY - posY)
}*/

func (rl *Instance) moveCursorByAdjust(adjust int) {
	oldX, oldY := lineWrapPos(rl.promptLen, rl.pos, rl.termWidth)

	rl.pos += adjust

	if rl.pos < 0 {
		rl.pos = 0
	}
	if rl.pos > len(rl.line) {
		rl.pos = len(rl.line)
	}

	if rl.modeViMode != vimInsert && rl.pos == len(rl.line) {
		rl.pos--
	}

	newX, newY := lineWrapPos(rl.promptLen, rl.pos, rl.termWidth)

	y := newY - oldY
	switch {
	case y < 0:
		moveCursorUp(-y)
	case y > 0:
		moveCursorDown(y)
	}

	x := newX - oldX
	switch {
	case x < 0:
		moveCursorBackwards(-x)
	case x > 0:
		moveCursorForwards(x)
	}
}
