package ansible

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/GearBoxLab/core/command"
	"github.com/GearBoxLab/core/printer"

	"github.com/symfony-cli/terminal"
)

type Ansible struct {
	commandFactory command.Factory
}

func New(commandFactory command.Factory) *Ansible {
	return &Ansible{
		commandFactory: commandFactory,
	}
}

func (ansible *Ansible) Install(osName, sudoPassword string) (err error) {
	var installed bool

	if installed, err = ansible.isInstalled(); nil != err {
		return err
	}

	if false == installed {
		switch osName {
		case "oracle-linux":
			commands := []*command.Command{
				ansible.commandFactory.NewSudoCommand(sudoPassword, "dnf", "check-update", "-y"),
				ansible.commandFactory.NewSudoCommand(sudoPassword, "dnf", "upgrade", "-y"),
				ansible.commandFactory.NewSudoCommand(sudoPassword, "dnf", "install", "-y", "epel-release"),
				ansible.commandFactory.NewSudoCommand(sudoPassword, "dnf", "install", "-y", "ansible"),
			}

			printer.Printf("\n<comment>$ %s</comment>\n", commands[0].String())
			if err = commands[0].Run(); err != nil {
				var exitError *exec.ExitError
				if errors.As(err, &exitError) && exitError.ExitCode() == 100 {
					printer.Printf("\n<comment>$ %s</comment>\n", commands[1].String())
					if err = commands[1].Run(); err != nil {
						return err
					}
				} else {
					return err
				}
			}

			printer.Printf("\n<comment>$ %s</comment>\n", commands[2].String())
			if err = commands[2].Run(); nil != err {
				return err
			}

			printer.Printf("\n<comment>$ %s</comment>\n", commands[3].String())
			if err = commands[3].Run(); nil != err {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("unsupported os: %q", osName))
		}
	}

	return nil
}

func (ansible *Ansible) RunAnsiblePlaybook(playbookFilePath, variableFilePath, sudoPassword string, args ...string) (err error) {
	cmd := ansible.commandFactory.NewCommand(
		"HISTSIZE=0",
		"ansible-playbook",
		playbookFilePath,
		"--extra-vars", "@"+variableFilePath,
		"--extra-vars", "ansible_become_password="+sudoPassword,
	)
	cmd.SetSecretArguments(6)

	if terminal.IsVerbose() {
		cmd.AddArguments("-" + strings.Repeat("v", terminal.GetLogLevel()-1))
	}

	if len(args) > 0 {
		cmd.AddArguments(args...)
	}

	printer.Printf("\n<comment>$ %s</comment>\n", cmd.String())
	if err = cmd.Run(); nil != err {
		return err
	}

	return nil
}

func (ansible *Ansible) isInstalled() (installed bool, err error) {
	var path string
	var realPath string

	if path, err = ansible.commandFactory.NewCommand("which", "ansible").Output(); nil != err && "exit status 1" != err.Error() {
		return false, err
	}
	path = strings.TrimSpace(path)

	if "" != path {
		if realPath, err = ansible.commandFactory.NewCommand("ls", path).Output(); nil != err {
			return false, err
		}

		realPath = strings.TrimSpace(realPath)

		return path == realPath, nil
	}

	return false, nil
}
