package utils

import (
	"fmt"
	"math/rand"
)

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomIP() string {
	ip := ""
	for i := 0; i < 4; i++ {
		ip += fmt.Sprint(rand.Intn(256))
		if i < 3 {
			ip += "."
		}
	}
	return ip
}

func RandomPort() int {
	return RandomInt(1024, 65535)
}
