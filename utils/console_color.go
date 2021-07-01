package utils

import "fmt"

const (
	colorRed  = "\u001B[31m"
	colorBlue = "\u001B[34m"
	reset     = "\u001B[0m"
)

func PrintlnWarning(text string) {
	Println(colorRed, text)
}

func PrintInfo(text string) {
	Println(colorBlue, text)
}

func Println(color string, text string) {
	fmt.Println(color + text + reset)
}
