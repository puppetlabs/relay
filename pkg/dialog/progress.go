package dialog

import (
	"fmt"
	"io"
	"time"

	"github.com/fatih/color"
)

var chars = []string{
	"⠋",
	"⠙",
	"⠹",
	"⠸",
	"⠼",
	"⠴",
	"⠦",
	"⠧",
	"⠇",
	"⠏",
}

var ProgressFrameDuration = 250 * time.Millisecond

// Progress is a text-based progress indicator that has a little animation
// associated with it. It's really cool.
type Progress struct {
	c *color.Color

	w   io.Writer
	msg string
	pos int

	// this is the channel that our ticks will come in on from the time package.
	// Once it's closed we know we can stop.
	ticks *time.Ticker
	done  chan bool
}

func (p *Progress) setNextPos() {
	p.pos += 1

	if p.pos >= len(chars) {
		p.pos = 0
	}
}

func (p *Progress) doRenderFrame() {
	fmt.Fprintf(p.w, "%s ", p.msg)
	p.c.Fprintf(p.w, "%s\r", chars[p.pos])
	p.setNextPos()
}

func (p *Progress) doRender() {
	for {
		select {
		case <-p.done:
			return
		case <-p.ticks.C:
			p.doRenderFrame()
		}
	}
}

func (p *Progress) Start() {
	p.ticks = time.NewTicker(ProgressFrameDuration)
	go p.doRender()
}

func (p *Progress) Complete() {
	p.done <- true

	// note that even though the loop has stopped we want to stop the underlying
	// ticker. this tells the time package that it can (if need be) be cleaned
	// up.
	p.ticks.Stop()

	fmt.Fprintf(p.w, "\r%s DONE!\n", p.msg)
}

func NewProgress(w io.Writer, msg string) *Progress {
	c := color.New(color.FgHiMagenta)

	return &Progress{
		c:    c,
		done: make(chan bool, 0),
		w:    w,
		msg:  msg,
	}
}
