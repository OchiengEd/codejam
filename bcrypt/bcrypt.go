package main

import (
  "fmt"
  "golang.org/x/crypto/bcrypt"
)

func main() {
  s := "pass124"
  hashedPass, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
  if err != nil {
    panic(err)
  }
  fmt.Println(string(hashedPass))
}
