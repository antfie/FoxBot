package crypto

import (
	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/blake2b"
)

func HashDataToString(data []byte) (string, error) {
	hash, err := blake2b.New512([]byte{})

	if err != nil {
		return "", err
	}

	hash.Write(data)

	return base58.Encode(hash.Sum(nil)), err
}
