package wsl

import (
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/GearBoxLab/core/process"
)

func DefaultWslDistribution() (distribution string, err error) {
	var result string

	p := process.New("wsl", "--list", "--verbose")
	regex := regexp.MustCompile(`^\* (\S+)\s+(Running|Stopped)\s+\d`)

	if result, err = p.Output(); nil != err {
		return "", err
	}

	for _, line := range strings.Split(result, "\n") {
		if matches := regex.FindStringSubmatch(line); len(matches) > 0 {
			return matches[1], nil
		}
	}

	return "", errors.New("cannot find default distribution name")
}

func ConvertToLinuxPath(distribution, windowsPath string) (linuxPath string, err error) {
	p := process.NewWslProcess(distribution, "wslpath", "-a", filepath.ToSlash(windowsPath))

	if linuxPath, err = p.Output(); nil != err {
		return "", err
	}

	return strings.TrimSpace(linuxPath), nil
}

func EnableSystemd(distribution string) (err error) {
	var content string

	args := []string{"[", "-f", "/etc/wsl.conf", "]", "&&", "cat", "/etc/wsl.conf"}
	if content, err = process.NewWslSudoProcess(distribution, args...).Output(); err != nil {
		return err
	}

	content = strings.TrimSpace(content)
	modified := false

	if "" == content {
		if err = createWslConfFile(distribution); err != nil {
			return err
		}
		modified = true
	} else {
		if modified, err = updateWslConfFile(distribution, content); err != nil {
			return err
		}
	}

	if modified {
		var terminateResult string

		if terminateResult, err = process.New("wsl", "--terminate", distribution).Output(); err != nil {
			if "The operation completed successfully." != strings.TrimSpace(terminateResult) {
				return err
			}
		}
	}

	return nil
}

func createWslConfFile(distribution string) (err error) {
	processes := []*process.Process{
		process.NewWslSudoProcess(distribution, "bash", "-c", "echo [boot] > /etc/wsl.conf"),
		process.NewWslSudoProcess(distribution, "bash", "-c", "echo systemd=true >> /etc/wsl.conf"),
		process.NewWslSudoProcess(distribution, "bash", "-c", "echo >> /etc/wsl.conf"),
	}
	for _, proc := range processes {
		if _, err = proc.Run(); err != nil {
			return err
		}
	}

	return nil
}

func updateWslConfFile(distribution, content string) (modified bool, err error) {
	var cfg *ini.File

	if cfg, err = ini.Load([]byte(content)); err != nil {
		return modified, err
	}

	if false == cfg.HasSection("boot") {
		if _, err = cfg.NewSection("boot"); err != nil {
			return modified, err
		}
	}

	if false == cfg.Section("boot").HasKey("systemd") {
		if _, err = cfg.Section("boot").NewKey("systemd", "true"); err != nil {
			return modified, err
		}
		modified = true
	} else if "true" != cfg.Section("boot").Key("systemd").Value() {
		cfg.Section("boot").Key("systemd").SetValue("true")
		modified = true
	}

	if modified {
		buff := strings.Builder{}
		if _, err = cfg.WriteTo(&buff); err != nil {
			return modified, err
		}

		lines := strings.Split(buff.String(), "\r\n")
		operatorCount := 1

		for index, line := range lines {
			if index > 0 {
				operatorCount = 2
			}

			cmd := fmt.Sprintf("echo %q %s /etc/wsl.conf", line, strings.Repeat(">", operatorCount))
			if _, err = process.NewWslSudoProcess(distribution, "bash", "-c", cmd).Run(); err != nil {
				return modified, err
			}
		}
	}

	return modified, nil
}
