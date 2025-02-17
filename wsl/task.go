package wsl

import (
	"errors"
	"fmt"
	"log"
	"os/user"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func CreateTaskForStartSystemd(distribution string) (err error) {
	var taskService *ole.IDispatch
	var rootFolder *ole.IDispatch
	taskFolderName := "WSL"
	taskName := fmt.Sprintf("Start %s systemd", distribution)
	taskPath := fmt.Sprintf(`\%s\%s`, taskFolderName, taskName)
	triggerTypeLogon := uint(9)
	actionTypeExecutable := uint(0)
	taskFlags := 6     // TASK_CREATE_OR_UPDATE
	taskLogonType := 3 // TASK_LOGON_INTERACTIVE_TOKEN

	if taskService, err = getTaskService(); err != nil {
		return err
	}
	defer taskService.Release()

	if rootFolder, err = getRootFolder(taskService); err != nil {
		return err
	}
	defer rootFolder.Release()

	if !hasSubFolder(rootFolder, taskFolderName) {
		if err = createSubFolder(rootFolder, taskFolderName); err != nil {
			log.Fatal(err)
		}
	}

	if hasTask(rootFolder, taskPath) {
		if _, err = oleutil.CallMethod(rootFolder, "DeleteTask", taskPath, 0); err != nil {
			return err
		}
	}

	var username string
	if username, err = getUsername(); err != nil {
		return err
	}

	task := oleutil.MustCallMethod(taskService, "NewTask", 0).ToIDispatch()
	defer task.Release()

	registrationInfo := oleutil.MustGetProperty(task, "RegistrationInfo").ToIDispatch()
	defer registrationInfo.Release()
	oleutil.MustPutProperty(registrationInfo, "Author", username)

	settings := oleutil.MustGetProperty(task, "Settings").ToIDispatch()
	defer settings.Release()
	oleutil.MustPutProperty(settings, "StartWhenAvailable", true)

	triggers := oleutil.MustGetProperty(task, "Triggers").ToIDispatch()
	defer triggers.Release()
	trigger := oleutil.MustCallMethod(triggers, "Create", triggerTypeLogon).ToIDispatch()
	defer trigger.Release()
	oleutil.MustPutProperty(trigger, "UserId", username)

	actions := oleutil.MustGetProperty(task, "Actions").ToIDispatch()
	defer actions.Release()
	action := oleutil.MustCallMethod(actions, "Create", actionTypeExecutable).ToIDispatch()
	defer action.Release()
	oleutil.MustPutProperty(action, "Path", "wsl")
	oleutil.MustPutProperty(action, "Arguments", fmt.Sprintf(`-d %s echo "Starting systemd..."`, distribution))

	_, err = oleutil.CallMethod(
		rootFolder,
		"RegisterTaskDefinition",
		taskPath,
		task,
		taskFlags,
		username,
		nil,
		taskLogonType,
	)
	if err != nil {
		return err
	}

	return nil
}

func getTaskService() (service *ole.IDispatch, err error) {
	if err = ole.CoInitialize(uintptr(0)); err != nil {
		code := err.(*ole.OleError).Code()
		if code != ole.S_OK && code != 0x00000001 {
			return nil, err
		}
	}

	var schedulerClassId *ole.GUID

	if schedulerClassId, err = ole.ClassIDFrom("Schedule.Service.1"); err != nil {
		ole.CoUninitialize()
		return nil, err
	}

	var scheduler *ole.IUnknown

	if scheduler, err = ole.CreateInstance(schedulerClassId, nil); err != nil {
		ole.CoUninitialize()
		return nil, err
	}

	if scheduler == nil {
		ole.CoUninitialize()
		return nil, errors.New("could not create ITaskService instance")
	}

	defer scheduler.Release()

	taskService := scheduler.MustQueryInterface(ole.IID_IDispatch)

	oleutil.MustCallMethod(taskService, "Connect", "", "", "", "")

	return taskService, nil
}

func getRootFolder(scheduleService *ole.IDispatch) (rootFolder *ole.IDispatch, err error) {
	var result *ole.VARIANT

	if result, err = oleutil.CallMethod(scheduleService, "GetFolder", `\`); err != nil {
		return nil, err
	}

	return result.ToIDispatch(), nil
}

func hasSubFolder(folder *ole.IDispatch, path string) bool {
	if _, err := oleutil.CallMethod(folder, "GetFolder", path); err != nil {
		return false
	}

	return true
}

func createSubFolder(folder *ole.IDispatch, path string) (err error) {
	if _, err = oleutil.CallMethod(folder, "CreateFolder", path, ""); err != nil {
		return err
	}

	return nil
}

func hasTask(folder *ole.IDispatch, path string) bool {
	if _, err := oleutil.CallMethod(folder, "GetTask", path); err != nil {
		return false
	}

	return true
}

func getUsername() (username string, err error) {
	var currentUser *user.User

	if currentUser, err = user.Current(); err != nil {
		return "", err
	}

	parts := strings.Split(currentUser.Username, `\`)

	if len(parts) > 0 {
		if len(parts) == 1 {
			return parts[0], nil
		} else {
			return parts[1], nil
		}
	}

	return username, err
}
