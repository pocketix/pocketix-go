package interfaces

type Command interface {
	Execute() error
	HasBlock() bool
}
