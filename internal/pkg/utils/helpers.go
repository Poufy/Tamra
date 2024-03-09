package utils

import (
	"fmt"
	"math/rand"
)

func GenerateCode() string {
	return fmt.Sprint(rand.Intn(900000) + 100000)
}
