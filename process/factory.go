package process

type Factory interface {
	NewProcess(arguments ...string) *Process
	NewSudoProcess(sudoPassword string, arguments ...string) *Process
}

type DefaultFactory struct {
}

func NewFactory() *DefaultFactory {
	return &DefaultFactory{}
}

func (f *DefaultFactory) NewProcess(arguments ...string) *Process {
	return New(arguments...)
}

func (f *DefaultFactory) NewSudoProcess(sudoPassword string, arguments ...string) *Process {
	return NewSudoProcess(sudoPassword, arguments...)
}

type WslFactory struct {
	distribution string
}

func NewWslFactory(distribution string) *WslFactory {
	return &WslFactory{distribution: distribution}
}

func (w *WslFactory) NewProcess(arguments ...string) *Process {
	return NewWslProcess(w.distribution, arguments...)
}

func (w *WslFactory) NewSudoProcess(sudoPassword string, arguments ...string) *Process {
	return NewWslSudoProcess(w.distribution, arguments...)
}
