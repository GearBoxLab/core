package ansible

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/GearBoxLab/core/process"
	"github.com/GearBoxLab/core/terminal"
)

type Ansible struct {
	processFactory process.Factory
}

func New(processFactory process.Factory) *Ansible {
	return &Ansible{
		processFactory: processFactory,
	}
}

func (i *Ansible) Install(osName, sudoPassword string) (err error) {
	var installed bool

	if installed, err = i.isInstalled(); nil != err {
		return err
	}

	if false == installed {
		switch osName {
		case "oracle-linux":
			processes := []*process.Process{
				i.processFactory.NewSudoProcess(sudoPassword, "dnf", "check-update", "-y"),
				i.processFactory.NewSudoProcess(sudoPassword, "dnf", "upgrade", "-y"),
				i.processFactory.NewSudoProcess(sudoPassword, "dnf", "install", "-y", "epel-release"),
				i.processFactory.NewSudoProcess(sudoPassword, "dnf", "install", "-y", "ansible"),
			}

			terminal.Printf("\n<comment>$ %s</comment>\n", processes[0].String())
			if _, err = processes[0].Run(); err != nil {
				var exitError *exec.ExitError
				if errors.As(err, &exitError) && exitError.ExitCode() == 100 {
					terminal.Printf("\n<comment>$ %s</comment>\n", processes[1].String())
					if _, err = processes[1].Run(); err != nil {
						return err
					}
				} else {
					return err
				}
			}

			terminal.Printf("\n<comment>$ %s</comment>\n", processes[2].String())
			if _, err = processes[2].Run(); nil != err {
				return err
			}

			terminal.Printf("\n<comment>$ %s</comment>\n", processes[3].String())
			if _, err = processes[3].Run(); nil != err {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("unsupported os: %q", osName))
		}
	}

	return nil
}

func (i *Ansible) RunAnsiblePlaybook(playbookFilePath, variableFilePath, sudoPassword string) (err error) {
	proc := i.processFactory.NewProcess(
		"HISTSIZE=0",
		"ansible-playbook",
		playbookFilePath,
		"--extra-vars", "@"+variableFilePath,
		"--extra-vars", "ansible_become_password="+sudoPassword,
	)
	proc.SetSecretArguments(6)

	if terminal.IsVerbose() {
		proc.AddArguments("-" + strings.Repeat("v", terminal.GetLogLevel()-1))
	}

	terminal.Printf("\n<comment>$ %s</comment>\n", proc.String())
	if _, err = proc.Run(); nil != err {
		return err
	}

	return nil
}

func (i *Ansible) isInstalled() (installed bool, err error) {
	var path string
	var realPath string

	if path, err = i.processFactory.NewProcess("which", "ansible").Output(); nil != err && "exit status 1" != err.Error() {
		return false, err
	}
	path = strings.TrimSpace(path)

	if "" != path {
		if realPath, err = i.processFactory.NewProcess("ls", path).Output(); nil != err {
			return false, err
		}

		realPath = strings.TrimSpace(realPath)

		return path == realPath, nil
	}

	return false, nil
}
