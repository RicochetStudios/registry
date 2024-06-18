package util

import "fmt"

const clientPrefix = "[client]"

func InfoMessage(m string) {
	InfoMessagef("%s\n", m)
}

func InfoMessagef(f string, a ...any) {
	fmt.Printf(clientPrefix+"[info] "+f, a...)
}

func WarningMessage(m string) {
	WarningMessagef("%s\n", m)
}

func WarningMessagef(f string, a ...any) {
	fmt.Printf(clientPrefix+"[warn] "+f, a...)
}

func ErrorMessage(m string) {
	ErrorMessagef("%s\n", m)
}

func ErrorMessagef(f string, a ...any) {
	fmt.Printf(clientPrefix+"[error] "+f, a...)
}
