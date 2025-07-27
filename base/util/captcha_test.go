package util

import (
	"os"
	"strings"
	"testing"

	"bbs-back/base/storage"
)

func TestCreateCaptcha(t *testing.T) {
	// Skip test if Redis is not available
	if storage.GetRedisPool() == nil {
		t.Skip("Redis is not available, skipping test")
	}

	// Try to ping Redis
	_, err := storage.GetRedisPool().Get("test_key")
	if err != nil && strings.Contains(err.Error(), "connection refused") {
		t.Skip("Redis connection refused, skipping test")
	}

	InitCaptcha()
	file, err := os.OpenFile("test.png", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()
	_, err = CreateCaptcha(file)
	if err != nil {
		t.Error(err)
		return
	}
}
