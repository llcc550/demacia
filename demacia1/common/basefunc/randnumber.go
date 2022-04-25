package basefunc

import (
	"math/rand"
	"strconv"
	"time"
)

func RandNumberString() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(int(rand.Uint32()))
}

func RandNumber() int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Uint32())
}
