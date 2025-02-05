package uac

import (
	"errors"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

// Prompt triggers Windows UAC elevation prompt.
func Prompt(messageFilePath string, maxWaitTime time.Duration, job func() error) error {
	return PromptWithExtraArguments(messageFilePath, maxWaitTime, []string{}, job)
}

// PromptWithExtraArguments triggers Windows UAC elevation prompt.
// The extraArguments will append to the original command arguments.
func PromptWithExtraArguments(messageFilePath string, maxWaitTime time.Duration, extraArguments []string, job func() error) error {
	if _, err := os.Stat(messageFilePath); nil == err {
		if err = os.Remove(messageFilePath); err != nil {
			return err
		}
	}

	if IsElevated() {
		err := job()
		message := []byte("")

		if nil != err {
			message = []byte(err.Error())
		}

		if writeErr := os.WriteFile(messageFilePath, message, 0644); writeErr != nil {
			return writeErr
		}

		return err
	}

	if err := doPrompt(extraArguments); err != nil {
		return err
	}

	spent := 0 * time.Millisecond
	step := 100 * time.Millisecond

	for {
		time.Sleep(step)
		spent += step

		if _, err := os.Stat(messageFilePath); errors.Is(err, os.ErrNotExist) {
			if spent >= maxWaitTime {
				break
			}
		} else {
			time.Sleep(step)

			if message, readErr := os.ReadFile(messageFilePath); readErr != nil {
				return readErr
			} else if len(message) > 0 {
				err = errors.New(string(message))
			}

			if removeErr := os.Remove(messageFilePath); removeErr != nil {
				return removeErr
			}

			if err != nil {
				return err
			}

			break
		}
	}

	return nil
}

func IsElevated() bool {
	return windows.GetCurrentProcessToken().IsElevated()
}

func doPrompt(extraArguments []string) error {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(append(os.Args[1:], extraArguments...), " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 // SW_NORMAL

	return windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
}
