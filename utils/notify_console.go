package utils

import (
	"github.com/fatih/color"
	"log"
)

func NotifyConsole(message string) {
	log.Print(color.CyanString(message))
}

func NotifyGoodConsole(message string) {
	log.Print(color.GreenString(message))
}

func NotifyBadConsole(message string) {
	log.Print(color.HiRedString(message))
}
