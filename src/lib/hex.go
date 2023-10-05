package lib

import "encoding/hex"

func GetPasswordFromHexString(hexPassword string) ([]byte, error) {
	byteData, err := hex.DecodeString(hexPassword)
	return byteData, err
}
