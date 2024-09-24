package utils

import (
	"os/exec"
	"path"
)

func NotifyAudio(message string) {
	play("hero_simple-celebration-01")
}

func NotifyGoodAudio(message string) {
	play("notification_decorative-01")
}

func NotifyBadAudio(message string) {
	play("alert_error-01")
}

func play(sound string) {
	// .wav files from https://m2.material.io/design/sound/sound-resources.html - "Material Design Sound Resources"
	cmd := exec.Command("/usr/bin/afplay", path.Join("alerts", sound+".wav")) // #nosec G204

	err := cmd.Run()

	if err != nil {
		// Sometimes this will error, no big deal.
		//log.Panic(err)
	}
}
