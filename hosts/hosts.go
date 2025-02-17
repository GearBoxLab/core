package hosts

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
)

const startMessageTemplate = `##>>> %s ## START >>>##`
const endMessageTemplate = `##<<< %s ## END   <<<##`

func IsHostsChanged(comment, newHostsFilePath, systemHostsFilePath string) (changed bool, err error) {
	var oldHostsBuff bytes.Buffer
	var newHosts []byte
	startMessage := fmt.Sprintf(startMessageTemplate, comment)
	endMessage := fmt.Sprintf(endMessageTemplate, comment)

	if _, oldHostsBuff, err = loadHostsFile(systemHostsFilePath, startMessage, endMessage); err != nil {
		return changed, err
	}
	if newHosts, err = os.ReadFile(newHostsFilePath); err != nil {
		return changed, err
	}

	oldHosts := bytes.TrimSpace(oldHostsBuff.Bytes())
	newHosts = bytes.TrimSpace(newHosts)
	newHosts = bytes.ReplaceAll(newHosts, []byte("\r\n"), []byte("\n"))

	return !bytes.Equal(oldHosts, newHosts), nil
}

func UpdateHostsFile(comment, newHostsFilePath, systemHostsFilePath string) (err error) {
	var newHosts []byte
	var buff bytes.Buffer
	startMessage := fmt.Sprintf(startMessageTemplate, comment)
	endMessage := fmt.Sprintf(endMessageTemplate, comment)

	if buff, _, err = loadHostsFile(systemHostsFilePath, startMessage, endMessage); err != nil {
		return err
	}

	if newHosts, err = os.ReadFile(newHostsFilePath); err != nil {
		return err
	}

	newHosts = bytes.TrimSpace(newHosts)

	if len(newHosts) > 0 {
		buff.WriteString(startMessage)
		buff.WriteRune('\n')
		buff.Write(newHosts)
		buff.WriteRune('\n')
		buff.WriteString(endMessage)
		buff.WriteRune('\n')
	}

	return os.WriteFile(systemHostsFilePath, buff.Bytes(), 0644)
}

func loadHostsFile(filepath, startMessage, endMessage string) (buff bytes.Buffer, oldBuff bytes.Buffer, err error) {
	if _, err = os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return buff, oldBuff, fmt.Errorf("the hosts file %q is not found", filepath)
	}

	var file *os.File
	if file, err = os.Open(filepath); err != nil {
		return buff, oldBuff, err
	}

	hostsFileScanner := bufio.NewScanner(file)
	regexpStart := regexp.MustCompile("^" + regexp.QuoteMeta(startMessage))
	regexpEnd := regexp.MustCompile("^" + regexp.QuoteMeta(endMessage))
	hasStart := false
	hasEnd := false

	hostsFileScanner.Split(bufio.ScanLines)

	for hostsFileScanner.Scan() {
		line := hostsFileScanner.Text()

		if regexpStart.MatchString(line) {
			hasStart = true
			continue
		}
		if regexpEnd.MatchString(line) {
			hasEnd = true
			continue
		}

		if (!hasStart && !hasEnd) || (hasStart && hasEnd) {
			buff.WriteString(line)
			buff.WriteRune('\n')
		} else {
			oldBuff.WriteString(line)
			oldBuff.WriteRune('\n')
		}
	}

	if err = file.Close(); err != nil {
		return buff, oldBuff, err
	}

	return buff, oldBuff, nil
}
