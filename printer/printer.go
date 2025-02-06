package printer

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

func Eprint(a ...interface{}) {
	if _, err := terminal.Eprint(a...); err != nil {
		log.Fatal(err)
	}
}

func Eprintf(format string, a ...interface{}) {
	if _, err := terminal.Eprintf(format, a...); err != nil {
		log.Fatal(err)
	}
}

func Eprintln(a ...interface{}) {
	if _, err := terminal.Eprintln(a...); err != nil {
		log.Fatal(err)
	}
}

func Eprintfln(format string, a ...interface{}) {
	if _, err := terminal.Eprintfln(format, a...); err != nil {
		log.Fatal(err)
	}
}
