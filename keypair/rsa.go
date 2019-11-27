package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"encoding/asn1"
	"fmt"
	"os"
)

var bitSize = 4096

func newRSAKeyPair() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		fmt.Println(err)
	}

	return privateKey
}

func save(filename string, content []byte) {
	certKeyName, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file ", err.Error())
	}
	defer certKeyName.Close()
	certKeyName.Write(content)
}

func main() {
  privateRSAKey := newRSAKeyPair()

	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateRSAKey),
		},
	)
	save("tls.key", privateKeyPEM)

	// RSA PUBLIC KEY type
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(&privateRSAKey.PublicKey),
		},
	)
	save("tls.crt", publicKeyPEM)

	// PUBLIC KEY type
	publicKeyBytes, err := asn1.Marshal(privateRSAKey.PublicKey)
	if err  != nil {
		fmt.Println("Error marshaling cert1 ", err.Error())
	}

	publicKey1PEM := pem.EncodeToMemory(
		&pem.Block{
			Type: "PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)
	save("tls1.crt", publicKey1PEM)
 }
