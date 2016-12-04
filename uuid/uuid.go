package uuid

import (
	"crypto/rand"
	"fmt"
)

type UUID struct {
	value string
}

func NewUUID() UUID {
	u := make([]byte, 16)
	_, err := rand.Read(u[:16])
	if err != nil {
		panic(err)
	}
	u[8] = (u[8] | 0x80) & 0xBf
	u[6] = (u[6] | 0x40) & 0x4f
	return UUID{value: fmt.Sprintf("%x-%x-%x-%x-%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])}
}
