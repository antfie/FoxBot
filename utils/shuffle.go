package utils

import (
	"crypto/rand"
	"log"
	"math/big"
)

func ShuffleStringArray(slice []string) {
	for i := len(slice) - 1; i > 0; i-- {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))

		if err != nil {
			log.Fatal(err)
		}

		j := int(nBig.Int64())

		slice[i], slice[j] = slice[j], slice[i]
	}
}
