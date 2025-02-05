package terminal

import (
	"log"

	"github.com/symfony-cli/terminal"
)

func Print(a ...interface{}) {
	if _, err := terminal.Print(a...); err != nil {
		log.Fatal(err)
	}
}

func Printf(format string, a ...interface{}) {
	if _, err := terminal.Printf(format, a...); err != nil {
		log.Fatal(err)
	}
}

func Println(a ...interface{}) {
	if _, err := terminal.Println(a...); err != nil {
		log.Fatal(err)
	}
}

func Printfln(format string, a ...interface{}) {
	if _, err := terminal.Printfln(format+"\n", a...); err != nil {
		log.Fatal(err)
	}
}

func IsVerbose() bool {
	return terminal.IsVerbose()
}

func GetLogLevel() int {
	return terminal.GetLogLevel()
}

func FormatBlockMessage(format string, msg string) string {
	return terminal.FormatBlockMessage(format, msg)
}
