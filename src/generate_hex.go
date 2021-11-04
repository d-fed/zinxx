package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func main() {
	val, _ := randomHex(20)
	fmt.Println(val)
}
