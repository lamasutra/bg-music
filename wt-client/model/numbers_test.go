package model

import (
	"testing"
)

func TestNarrate(t *testing.T) {
	num := Heading(5)
	str := num.Narrate()

	if str[0] != "5" {
		t.Error("should be 5")
	}

	num = Heading(11)
	str = num.Narrate()
	if str[0] != "11" {
		t.Error("should be 11")
	}

	num = Heading(12)
	str = num.Narrate()
	if str[0] != "12" {
		t.Error("should be 12")
	}

	num = Heading(13)
	str = num.Narrate()
	if str[0] != "13" {
		t.Error("should be 13")
	}

	num = Heading(14)
	str = num.Narrate()
	if str[0] != "14" {
		t.Error("should be 14")
	}

	num = Heading(15)
	str = num.Narrate()
	if str[0] != "15" {
		t.Error("should be 15")
	}

	num = Heading(20)
	str = num.Narrate()
	if str[0] != "20" {
		t.Error("should be 20")
	}

	num = Heading(21)
	str = num.Narrate()
	if str[0] != "20" || str[1] != "1" {
		t.Error("should be 20 and 1")
	}

	num = Heading(315)
	str = num.Narrate()
	if str[0] != "3" || str[1] != "100" || str[2] != "15" {
		t.Error("should be 3, 100 and 15")
	}

}
