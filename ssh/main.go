package main

import (
  "crypto/rsa"
  "crypto/rand"
  "crypto/x509"
  "encoding/pem"
  "golang.org/x/crypto/ssh"
  "fmt"
  "crypto/dsa"
  "crypto/ecdsa"
  "crypto/elliptic"
  "encoding/asn1"
  "crypto/ed25519"
)

const (
  bitSize = 4096
)

func generateRSAKeyPair() ([]byte, []byte, error) {
  privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
  if err != nil {
    return nil, nil, err
  }

  block := &pem.Block{
    Type: "RSA PRIVATE KEY",
    Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
  }

  // Generate SSH public key
  publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
  if err != nil {
    return nil, nil, err
  }

  publicKeyBytes := ssh.MarshalAuthorizedKey(publicKey)

  return pem.EncodeToMemory(block), publicKeyBytes, nil
}


func generateECDSAKeyPair() ([]byte, []byte, error) {
  privatekey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	keyBytes, err := x509.MarshalECPrivateKey(privatekey)
	if err != nil {
		return nil, nil, err
	}

	block := &pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: keyBytes,
	}

	publicKey, err := ssh.NewPublicKey(&privatekey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return pem.EncodeToMemory(block), ssh.MarshalAuthorizedKey(publicKey), nil
}


func generateDSAKeyPair()  ([]byte, []byte, error) {
  params := new(dsa.Parameters)

	if err := dsa.GenerateParameters(params, rand.Reader, dsa.L1024N160); err != nil {
		return nil, nil, err
	}

	privateKey := new(dsa.PrivateKey)
	privateKey.PublicKey.Parameters = *params

	if err := dsa.GenerateKey(privateKey, rand.Reader); err != nil {
		return nil, nil, err
	}

	asnBytes, err := asn1.Marshal(*privateKey)
	if err != nil {
		return nil, nil, err
	}

	block := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: asnBytes,
	}

	// publicKey := privateKey.PublicKey
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return pem.EncodeToMemory(block), ssh.MarshalAuthorizedKey(publicKey), nil
}

func generateED25519KeyPair() ([]byte, []byte, error) {
  public, private, err := ed25519.GenerateKey(rand.Reader)
  if err != nil {
    return nil, nil, err
  }

  asnBytes, err := asn1.Marshal(private)
	if err != nil {
		return nil, nil, err
	}

  privatePEM := pem.Block{
    Type: "OPENSSH PRIVATE KEY",
    Bytes: asnBytes,
  }

  publicKey, err := ssh.NewPublicKey(public)
  if err != nil {
		return nil, nil, err
	}

  return pem.EncodeToMemory(&privatePEM), ssh.MarshalAuthorizedKey(publicKey), nil
}

func main() {
  fmt.Println("Hello world")
  key, pub, err := generateRSAKeyPair()
  if err != nil {
    fmt.Printf("%+v", err)
  }
  fmt.Println(string(key))
  fmt.Println(string(pub))

}
