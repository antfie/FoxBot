package utils

import (
	"github.com/fatih/color"
	"log"
)

func Notify(message string) {
	log.Print(color.CyanString(message))
	//NotifySlack(message)
	//play("hero_simple-celebration-01")
}

func NotifyGood(message string) {
	log.Print(color.GreenString(message))
	//NotifySlack(message)
	// play("notification_decorative-01")
}

func NotifyBad(message string) {
	log.Print(color.HiRedString(message))
	//NotifySlack(message)
	// play("alert_error-01")
}
