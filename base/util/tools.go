package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GetRandomString(l int) string {
	resources := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	result := make([]byte, 0, l)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, resources[r.Intn(len(resources))])
	}
	return string(result)
}

func GetEmailKey(account string) string {
	return fmt.Sprintf("Email:%s", account)
}
