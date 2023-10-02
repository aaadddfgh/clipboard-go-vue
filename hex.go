package main

import "encoding/hex"

func getPasswordFromHexString(hexPassword string) ([]byte, error) {
	byteData, err := hex.DecodeString(hexPassword)
	return byteData, err
}
