package peon

type Job interface {
	Execute() error
}
