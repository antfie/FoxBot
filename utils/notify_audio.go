package utils

import (
	"os/exec"
	"path"
)

func play(sound string) {
	// .wav files from https://m2.material.io/design/sound/sound-resources.html - "Material Design Sound Resources"
	cmd := exec.Command("/usr/bin/afplay", path.Join("alerts", sound+".wav")) // #nosec G204

	err := cmd.Run()

	if err != nil {
		// Sometimes this will error, no big deal.
		//log.Panic(err)
	}
}
