package matcha

type Context struct {
	id        string
	hookIndex int
	store     *store
	throttler *throttler
}
