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

func TestSlalomEmail(t *testing.T) {
	bad := "not@email.com"
	good := "slkn@slalom.com"
	if validSlalomEmail(bad) == true {
		t.Errorf("%s is not a valid email.", bad)
	}

	if validSlalomEmail(good) == false {
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

func TestIsUUID(t *testing.T) {
	UUID := "017d9a04-7dd9-40a9-a4e1-2b0dacaa46db"
	notUUID := "bad"
	good := IsValidUUID(UUID)
	if !good {
		t.Errorf("expected: true got %t", good)
	}
	bad := IsValidUUID(notUUID)
	if bad {
		t.Errorf("expected: true got %t", bad)
	}
}
