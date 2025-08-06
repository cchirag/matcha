package core

type Component interface {
	Render(ctx *Ctx) Component
}
