package utils

import (
	"fmt"
	"testing"
)

func TestVal(t *testing.T) {
	phone := "1025906123"
	fmt.Println(VerifyMobileFormat(phone))
}

func TestCompareSimple(t *testing.T) {
	fmt.Println(CompareSimple("12", "11.9"))
}
