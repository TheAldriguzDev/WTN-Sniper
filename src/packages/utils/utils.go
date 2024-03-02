package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GetTimeStamp() string {
	currentTime := time.Now()
	timestamp := currentTime.UnixNano() / int64(time.Millisecond)
	cacheBuster := strconv.Itoa(int(timestamp))
	return cacheBuster
}

func GenNum(min, max int) int {
	// Imposta un seme diverso per ogni esecuzione
	rand.Seed(time.Now().UnixNano())

	// Genera un numero casuale compreso tra min e max
	return rand.Intn(max-min+1) + min
}
