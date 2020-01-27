package main

import (
  "fmt"
  "github.com/Nerzal/gocloak"
)

func main()  {
  client := gocloak.NewClient("http://172.17.0.2:8080")
  token, err := client.Login("demo-client", "07eb1b16-cf7a-44ab-9f16-8826d57feeed", "demo", "demo", "demo")
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  if err := client.Logout("demo-client", "07eb1b16-cf7a-44ab-9f16-8826d57feeed", "demo", token.RefreshToken); err != nil {
    fmt.Println(err.Error())
    return
  }

  fmt.Println("Successfully logged out user")
}
