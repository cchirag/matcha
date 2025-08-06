package core

type Ctx struct {
	States    []string
	Id        string
	ParentID  string
	HookIndex int
}
