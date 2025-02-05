package hosts

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
)

func UpdateHostsFile(comment, newHostsFilePath, systemHostsFilePath string) (err error) {
	var newHosts []byte
	var buff bytes.Buffer
	startMessage := fmt.Sprintf(`##>>> %s ## START >>>##`, comment)
	endMessage := fmt.Sprintf(`##<<< %s ## END   <<<##`, comment)

	if buff, err = loadHostsFile(systemHostsFilePath, startMessage, endMessage); err != nil {
		return err
	}

	if newHosts, err = os.ReadFile(newHostsFilePath); nil != err {
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

func loadHostsFile(filepath, startMessage, endMessage string) (buff bytes.Buffer, err error) {
	if _, err = os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return buff, fmt.Errorf("the hosts file %q is not found", filepath)
	}

	var file *os.File
	if file, err = os.Open(filepath); nil != err {
		return buff, err
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

		if (false == hasStart && false == hasEnd) || (true == hasStart && true == hasEnd) {
			buff.WriteString(line)
			buff.WriteRune('\n')
		}
	}

	if err = file.Close(); nil != err {
		return buff, err
	}

	return buff, nil
}
