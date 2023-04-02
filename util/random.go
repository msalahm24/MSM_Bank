package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	alphabets = "qwertyuiopasdfghjklzxcvbnm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generate random int between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Generate Random string
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabets)

	for i := 0; i < n; i++ {
		c := alphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwnerName() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 100000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(10))
}
