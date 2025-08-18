package matcha

type dimensions struct {
	Width, Height int
}

type Context struct {
	id         string
	hookIdx    int
	channels   *channels
	managers   *managers
	dimensions *dimensions
}

func (c *Context) Quit() {
	close(c.channels.quit)
}

func (c *Context) GetDimensions() *dimensions {
	return c.dimensions
}
