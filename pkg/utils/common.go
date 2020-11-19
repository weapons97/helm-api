package utils

import (
	"math/rand"
	"sync"
	"time"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyz")
	grand   *rand.Rand
)

func init() {
	s := rand.NewSource(time.Now().Unix())
	grand = rand.New(s)
}

func RandomString(n int) string {

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[grand.Intn(len(letters))]
	}

	return string(b)
}

var ConnectionWaiter = sync.WaitGroup{}
