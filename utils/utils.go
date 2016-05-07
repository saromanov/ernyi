package utils

import (
	"math/rand"
	"strconv"
)

// Some helpful tools

// GenRandomPort provides generation of the port
func GenRandomPort() string {
	return strconv.Itoa(randInt(10000, 65536))
}

func randInt(min int, max int) int {
    return min + rand.Intn(max-min)
}