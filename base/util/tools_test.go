package util

import (
	"fmt"
	"regexp"
	"testing"
)

func TestGetRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   int
		regex  string
	}{
		{
			name:   "Length 0",
			length: 0,
			want:   0,
			regex:  "^[0-9a-z]*$",
		},
		{
			name:   "Length 5",
			length: 5,
			want:   5,
			regex:  "^[0-9a-z]*$",
		},
		{
			name:   "Length 10",
			length: 10,
			want:   10,
			regex:  "^[0-9a-z]*$",
		},
		{
			name:   "Length 20",
			length: 20,
			want:   20,
			regex:  "^[0-9a-z]*$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRandomString(tt.length)
			
			// Check length
			if len(got) != tt.want {
				t.Errorf("GetRandomString() length = %v, want %v", len(got), tt.want)
			}
			
			// Check if the string matches the expected pattern
			match, err := regexp.MatchString(tt.regex, got)
			if err != nil {
				t.Errorf("Regex error: %v", err)
			}
			if !match {
				t.Errorf("GetRandomString() = %v, doesn't match pattern %v", got, tt.regex)
			}
		})
	}
	
	// Test randomness by generating multiple strings and checking they're different
	if GetRandomString(10) == GetRandomString(10) {
		t.Errorf("GetRandomString() generated identical strings, expected different random strings")
	}
}

func TestGetEmailKey(t *testing.T) {
	tests := []struct {
		name    string
		account string
		want    string
	}{
		{
			name:    "Empty account",
			account: "",
			want:    "Email:",
		},
		{
			name:    "Normal account",
			account: "test@example.com",
			want:    "Email:test@example.com",
		},
		{
			name:    "Account with special characters",
			account: "test.user+tag@example.com",
			want:    "Email:test.user+tag@example.com",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetEmailKey(tt.account)
			if got != tt.want {
				t.Errorf("GetEmailKey() = %v, want %v", got, tt.want)
			}
			
			// Also verify it matches the expected format
			expected := fmt.Sprintf("Email:%s", tt.account)
			if got != expected {
				t.Errorf("GetEmailKey() = %v, want %v", got, expected)
			}
		})
	}
}