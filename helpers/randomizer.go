package helpers

import (
	"math/rand"
	"time"
)

func RandomBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 0
}
