package main

import (
  "crypto/rand"
  "crypto/x509"
  "crypto/x509/pkix"
  "crypto/rsa"
  "encoding/pem"
  "time"
  "math/big"
  "log"
  "fmt"
  "os"
)

var bitSize = 2048

func createRSAKeyPair() *rsa.PrivateKey {
  RSAKeyPair, err := rsa.GenerateKey(rand.Reader, bitSize)
  if err != nil {
    log.Println("Error generating key ", err.Error())
  }

  return RSAKeyPair
}

func save(filename string, content []byte) {
  certKeyFile, err := os.Create(filename)
  if err != nil {
    fmt.Println("Error creating %s. ", filename, err.Error())
  }
  defer certKeyFile.Close()
  certKeyFile.Write(content)
}

func main() {
  privateKey := createRSAKeyPair()
  privateKeyPEM := pem.EncodeToMemory(
    &pem.Block {
      Type: "RSA PRIVATE KEY",
      Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    },
  )

  template := x509.Certificate {
    SerialNumber: big.NewInt(1),
    Subject: pkix.Name{
      Organization: []string{"Doorman"},
    },
    NotBefore: time.Now(),
    NotAfter: time.Now().Add(time.Hour * 24 * 180),
    KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
    ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
    BasicConstraintsValid: true,
    IsCA: false,
  }

  certificateDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
  if err != nil {
      log.Println("Error generating certificate: ", err.Error())
  }

  certificatePEM := pem.EncodeToMemory(
    &pem.Block{
        Type: "CERTIFICATE",
        Bytes: certificateDER,
    },
  )

  save("tls.key", privateKeyPEM)
  save("tls.crt", certificatePEM)
}
