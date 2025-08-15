package matcha

type Context struct {
	id       string
	channels *channels
	managers *managers
}

func (c *Context) Quit() {
	close(c.channels.quit)
}
