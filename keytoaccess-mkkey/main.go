package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	password := "my-secret-key"

	// Hash the password using SHA256
	h := sha256.New()
	h.Write([]byte(password))
	fmt.Println(hex.EncodeToString(h.Sum(nil)))
}
