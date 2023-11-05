package controller

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func genid(ext string) string {
	currentTime := time.Now().Format("20060102150405")
	randomString := randomString(12)
	return fmt.Sprintf("%s_%s%s", currentTime, randomString, ext)
}

func randomString(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteByte(chars[rand.Intn(len(chars))])
	}
	return result.String()
}