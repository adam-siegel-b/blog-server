package main

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func validEmail(email string) bool {
	stripped := strings.TrimSpace(email)
	_, err := mail.ParseAddress(stripped)
	return err == nil
}

func stripSketchyChars(value string) string {
	stripped := strings.TrimSpace(value)
	reg, _ := regexp.Compile("[^a-zA-Z0-9_-]+")
	processedString := reg.ReplaceAllString(stripped, "")
	return processedString
}

func validSlalomEmail(email string) bool {
	isEmail := validEmail(email)
	if isEmail {
		stripped := strings.TrimSpace(email)
		a := strings.Split(stripped, "@")
		if strings.ToLower(a[len(a)-1]) == "slalom.com" {
			return true
		}
	}
	return false
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
