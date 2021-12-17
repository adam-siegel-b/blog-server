package main

import (
	"net/mail"
	"regexp"
)

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func stripSketchyChars(value string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9_-]+")
	processedString := reg.ReplaceAllString(value, "")
	return processedString
}
