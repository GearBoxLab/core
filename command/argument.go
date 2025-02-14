package command

import (
	"fmt"
	"strings"
)

type argument struct {
	Value    string
	IsSecret bool
}

func (arg *argument) ToString(hideSecret bool) string {
	if arg.IsSecret && hideSecret {
		return "***secret***"
	}

	if strings.Contains(arg.Value, " ") || strings.Contains(arg.Value, "\"") {
		return fmt.Sprintf(`"%s"`, strings.ReplaceAll(arg.Value, `"`, `\"`))
	}

	return arg.Value
}
