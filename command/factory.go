package command

type Factory interface {
	NewCommand(arguments ...string) *Command
	NewSudoCommand(sudoPassword string, arguments ...string) *Command
}

type DefaultFactory struct {
}

func NewFactory() *DefaultFactory {
	return &DefaultFactory{}
}

func (f *DefaultFactory) NewCommand(arguments ...string) *Command {
	return New(arguments...)
}

func (f *DefaultFactory) NewSudoCommand(sudoPassword string, arguments ...string) *Command {
	return NewSudoCommand(sudoPassword, arguments...)
}

type WslFactory struct {
	distribution string
}

func NewWslFactory(distribution string) *WslFactory {
	return &WslFactory{distribution: distribution}
}

func (w *WslFactory) NewCommand(arguments ...string) *Command {
	return NewWslCommand(w.distribution, arguments...)
}

func (w *WslFactory) NewSudoCommand(sudoPassword string, arguments ...string) *Command {
	return NewWslSudoCommand(w.distribution, arguments...)
}
