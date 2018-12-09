package utils

import "github.com/fatih/color"

var Common, Red, Yellow *color.Color

func init() {
	Common = color.New(color.FgCyan)
	Red = color.New(color.FgRed)
	Yellow = color.New(color.FgGreen).Add(color.Bold)
}

func TwoColorPrintLn(common, yellow string) {
	Common.Print(common)
	Yellow.Println(yellow)
}
