package utils

import (
	"github.com/fatih/color"
	"log"
)

var Common, Red, Yellow, FgMagenta *color.Color

func init() {
	Common = color.New(color.FgCyan)
	Red = color.New(color.FgRed)
	Yellow = color.New(color.FgGreen).Add(color.Bold)
	FgMagenta = color.New(color.FgMagenta)
}

func TwoColorPrintLn(common, yellow string) {
	Common.Print(common)
	Yellow.Println(yellow)
}

func WhiteAndFgMagenta(white, fgMagenta string) {
	log.Print(white)
	FgMagenta.Println(fgMagenta)
}
