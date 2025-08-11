package matcha

import (
	"log"
)

type Context struct {
	id       string
	channels *channels
	managers *managers
}

func (c *Context) Quit() {
	log.Println("Quitting")
	c.channels.quit <- struct{}{}
}
