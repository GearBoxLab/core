package command

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	preArguments []*argument
	arguments    []*argument
}

func New(arguments ...string) *Command {
	p := &Command{
		preArguments: []*argument{},
		arguments:    []*argument{},
	}
	p.AddArguments(arguments...)

	return p
}

func NewSudoCommand(sudoPassword string, arguments ...string) *Command {
	cmd := New(arguments...).AddPreArguments("HISTSIZE=0", "echo", sudoPassword, "|", "sudo", "-S")
	cmd.SetSecretPreArguments(2)

	return cmd
}

func NewWslCommand(distribution string, arguments ...string) *Command {
	return New(arguments...).AddPreArguments("wsl", "-d", distribution)
}

func NewWslSudoCommand(distribution string, arguments ...string) *Command {
	return New(arguments...).AddPreArguments("wsl", "-d", distribution, "-u", "root")
}

func (p *Command) Run() (err error) {
	cmd := p.newExecCmd()
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (p *Command) Output() (out string, err error) {
	var result []byte

	if result, err = p.newExecCmd().CombinedOutput(); err != nil {
		return "", err
	}

	result = bytes.ReplaceAll(result, []byte{'\x00'}, []byte{})
	result = bytes.ReplaceAll(result, []byte{'\r'}, []byte{})

	return string(result), err
}

func (p *Command) String() string {
	return p.toString(true)
}

func (p *Command) StringWithSecret() string {
	return p.toString(false)
}

func (p *Command) AddPreArguments(arguments ...string) *Command {
	for _, arg := range arguments {
		p.preArguments = append(p.preArguments, &argument{Value: arg, IsSecret: false})
	}

	return p
}

func (p *Command) AddArguments(arguments ...string) *Command {
	for _, arg := range arguments {
		p.arguments = append(p.arguments, &argument{Value: arg, IsSecret: false})
	}

	return p
}

func (p *Command) SetSecretPreArguments(indexes ...int) *Command {
	lastIndex := len(p.preArguments) - 1

	for _, index := range indexes {
		if 0 <= index && index <= lastIndex {
			p.preArguments[index].IsSecret = true
		}
	}

	return p
}

func (p *Command) SetSecretArguments(indexes ...int) *Command {
	lastIndex := len(p.arguments) - 1

	for _, index := range indexes {
		if 0 <= index && index <= lastIndex {
			p.arguments[index].IsSecret = true
		}
	}

	return p
}

func (p *Command) SetNormalPreArguments(indexes ...int) *Command {
	lastIndex := len(p.preArguments) - 1

	for _, index := range indexes {
		if 0 <= index && index <= lastIndex {
			p.preArguments[index].IsSecret = false
		}
	}

	return p
}

func (p *Command) SetNormalArguments(indexes ...int) *Command {
	lastIndex := len(p.arguments) - 1

	for _, index := range indexes {
		if 0 <= index && index <= lastIndex {
			p.arguments[index].IsSecret = false
		}
	}

	return p
}

func (p *Command) newExecCmd() *exec.Cmd {
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

func (p *Command) toString(hideSecret bool) string {
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
