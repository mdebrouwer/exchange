package main

import (
	"encoding/hex"
	"fmt"

	"github.com/gorilla/securecookie"
)

func newKey() string {
	return hex.EncodeToString(securecookie.GenerateRandomKey(32))
}

func main() {
	fmt.Printf("Cookie Signing Key:    [%s]\n", newKey())
	fmt.Printf("Cookie Encryption key: [%s]\n", newKey())
}
