package main

import (
  "crypto/rsa"
  "crypto/rand"
  "crypto/x509"
  "encoding/pem"
  "golang.org/x/crypto/ssh"
  "log"
  "fmt"
)

var bitSize = 4096

func createRSAKeyPair() ([]byte, []byte) {
  privateRSAKey, err := rsa.GenerateKey(rand.Reader, bitSize)
  if err != nil {
    log.Println("Create key pair ", err.Error())
  }

  privateKeyPEM := pem.EncodeToMemory(
    &pem.Block{
      Type: "RSA PRIVATE KEY",
      Headers: nil,
      Bytes: x509.MarshalPKCS1PrivateKey(privateRSAKey),
    },
  )

  publicSSHKey, err := ssh.NewPublicKey(&privateRSAKey.PublicKey)
  if err != nil {
    log.Println("Error creating SSh public key ", err.Error())
  }

  publicKeyBytes := ssh.MarshalAuthorizedKey(publicSSHKey)

  return privateKeyPEM, publicKeyBytes
}

func main() {
  privateKey, publicKey := createRSAKeyPair()
  fmt.Println(string(privateKey))
  fmt.Println(string(publicKey))
}
