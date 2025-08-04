package utils

import (
	"github.com/fatih/color"
	"log"
)

func NotifyConsole(message string) {
	log.Print(color.CyanString(message))
}

func NotifyConsoleGood(message string) {
	log.Print(color.GreenString(message))
	go play("notification_decorative-01")
}

func NotifyConsoleBad(message string) {
	log.Print(color.HiRedString(message))
	go play("alert_error-01")
}
