package main

import "testing"

func TestValidEmail(t *testing.T) {
	bad := "notemail"
	good := "test@slalom.com"
	if validEmail(bad) == true {
		t.Errorf("%s is not a valid email.", bad)
	}

	if validEmail(good) == false {
		t.Errorf("%s is a valid email.", bad)
	}
}

func TestStripSketchyChars(t *testing.T) {
	normal := "normal"
	sketchy := "^*&^*very sketch&(*&(*"
	if stripSketchyChars(normal) != normal {
		t.Errorf("expected: %s got %s", normal, stripSketchyChars(normal))
	}
	if stripSketchyChars(sketchy) != "verysketch" {
		t.Errorf("expected: %s got %s", "verysketch", stripSketchyChars(sketchy))
	}
}
