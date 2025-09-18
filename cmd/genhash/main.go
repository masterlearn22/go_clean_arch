package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	h, _ := bcrypt.GenerateFromPassword([]byte("1234567890"), bcrypt.DefaultCost)
	fmt.Println(string(h))
}
