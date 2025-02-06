package process

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

type Process struct {
	preArguments []*argument
	arguments    []*argument
}

func New(arguments ...string) *Process {
	p := &Process{
		preArguments: []*argument{},
		arguments:    []*argument{},
	}
	p.AddArguments(arguments...)

	return p
}

func NewSudoProcess(sudoPassword string, arguments ...string) *Process {
	process := New(arguments...).AddPreArguments("HISTSIZE=0", "echo", sudoPassword, "|", "sudo", "-S")
	process.SetSecretPreArguments(2)

	return process
}

func NewWslProcess(distribution string, arguments ...string) *Process {
	return New(arguments...).AddPreArguments("wsl", "-d", distribution)
}

func NewWslSudoProcess(distribution string, arguments ...string) *Process {
	return New(arguments...).AddPreArguments("wsl", "-d", distribution, "-u", "root")
}

func (p *Process) Run() (out string, err error) {
	cmd := p.newCommand()
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return "", err
	}

	return "", nil
}

func (p *Process) Output() (out string, err error) {
	var result []byte

	if result, err = p.newCommand().CombinedOutput(); err != nil {
		return "", err
	}

	result = bytes.ReplaceAll(result, []byte{'\x00'}, []byte{})
	result = bytes.ReplaceAll(result, []byte{'\r'}, []byte{})

	return string(result), err
}

func (p *Process) String() string {
	return p.toString(true)
}

func (p *Process) StringWithSecret() string {
	return p.toString(false)
}

func (p *Process) AddPreArguments(arguments ...string) *Process {
	for _, arg := range arguments {
		p.preArguments = append(p.preArguments, &argument{Value: arg, IsSecret: false})
	}

	return p
}

func (p *Process) AddArguments(arguments ...string) *Process {
	for _, arg := range arguments {
		p.arguments = append(p.arguments, &argument{Value: arg, IsSecret: false})
	}

	return p
}

func (p *Process) SetSecretPreArguments(indexes ...int) *Process {
	lastIndex := len(p.preArguments) - 1

	for _, index := range indexes {
		if 0 <= index && index <= lastIndex {
			p.preArguments[index].IsSecret = true
		}
	}

	return p
}

func (p *Process) SetSecretArguments(indexes ...int) *Process {
	lastIndex := len(p.arguments) - 1

	for _, index := range indexes {
		if 0 <= index && index <= lastIndex {
			p.arguments[index].IsSecret = true
		}
	}

	return p
}

func (p *Process) SetNormalPreArguments(indexes ...int) *Process {
	lastIndex := len(p.preArguments) - 1

	for _, index := range indexes {
		if 0 <= index && index <= lastIndex {
			p.preArguments[index].IsSecret = false
		}
	}

	return p
}

func (p *Process) SetNormalArguments(indexes ...int) *Process {
	lastIndex := len(p.arguments) - 1

	for _, index := range indexes {
		if 0 <= index && index <= lastIndex {
			p.arguments[index].IsSecret = false
		}
	}

	return p
}

func (p *Process) newCommand() *exec.Cmd {
	args := make([]string, len(p.preArguments)+len(p.arguments))
	index := 0

	for _, arg := range p.preArguments {
		args[index] = arg.Value
		index++
	}

	for _, arg := range p.arguments {
		args[index] = arg.Value
		index++
	}

	return exec.Command(args[0], args[1:]...)
}

func (p *Process) toString(hideSecret bool) string {
	buffer := strings.Builder{}
	preArgumentLastIndex := len(p.preArguments) - 1
	argumentLastIndex := len(p.arguments) - 1

	for index, arg := range p.preArguments {
		buffer.WriteString(arg.ToString(hideSecret))

		if index != preArgumentLastIndex {
			buffer.WriteRune(' ')
		}
	}

	if preArgumentLastIndex > -1 {
		buffer.WriteRune(' ')
	}

	for index, arg := range p.arguments {
		buffer.WriteString(arg.ToString(hideSecret))

		if index != argumentLastIndex {
			buffer.WriteRune(' ')
		}
	}

	return buffer.String()
}
