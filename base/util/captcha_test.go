package util

import (
	"os"
	"testing"
)

func TestCreateCaptcha(t *testing.T) {
	Init()
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
