package style

type Style struct {
	bold      bool
	italics   bool
	underline bool
	fgcolor   Color
	bgcolor   Color
}

func (s *Style) Italics(value bool) *Style {
	s.italics = value
	return s
}

func (s *Style) GetItalice() bool {
	return s.italics
}

func (s *Style) Bold(value bool) *Style {
	s.bold = value
	return s
}

func (s *Style) GetBold() bool {
	return s.bold
}

func (s *Style) Underline(value bool) *Style {
	s.underline = value
	return s
}

func (s *Style) GetUnderline() bool {
	return s.underline
}

func (s *Style) FgColor(color Color) *Style {
	s.fgcolor = color
	return s
}

func (s *Style) GetFgColor() Color {
	return s.fgcolor
}

func (s *Style) BgColor(color Color) *Style {
	s.bgcolor = color
	return s
}

func (s *Style) GetBgColor() Color {
	return s.bgcolor
}
