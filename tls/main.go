package main

import (
  // "crypto/tls"
  "crypto/rand"
  "crypto/rsa"
  "encoding/pem"
  "crypto/x509"
  "fmt"
  "math/big"
  "time"
  "crypto/x509/pkix"
  "net"
)

var (
  keyBitSize = 4096
)

// PrivateKeyRSA returns a RSA private key
func PrivateKeyRSA() (*rsa.PrivateKey, error) {
  privateKey, err := rsa.GenerateKey(rand.Reader, keyBitSize)
  if err != nil {
    return nil, err
  }

  return privateKey, nil
}

// TLSCertificateAndKey returns certificate, privateKey and error
func TLSCertificateAndKey() ([]byte, []byte, error) {
  key, err := PrivateKeyRSA()
  if err != nil {
    return nil, nil, err
  }

  privateKeyPEM := pem.EncodeToMemory(
    &pem.Block{
      Type: "RSA PRIVATE KEY",
      Bytes: x509.MarshalPKCS1PrivateKey(key),
    },
  )

  cert := getCACertificate()
  ca := getClientCertificate()

 certificateDER, err := x509.CreateCertificate(rand.Reader, &cert, &ca, &key.PublicKey, key)
 if err != nil {
   return nil, nil, err
 }

 return privateKeyPEM,
  pem.EncodeToMemory(&pem.Block{
    Type: "CERTIFICATE",
    Bytes: certificateDER,
   }), nil
}

func main()  {
  fmt.Println("Hello world")
  key, cert, err := TLSCertificateAndKey()
  if err != nil {
    fmt.Println(err)
  }

  fmt.Println(string(key))
  fmt.Println(string(cert))
}

func getCACertificate() x509.Certificate {
  return x509.Certificate{
    SerialNumber: big.NewInt(1),
    Subject: pkix.Name{
      Organization: []string{"GitLab, Inc."},
      Country: []string{"US"},
    },
    NotBefore: time.Now(),
    NotAfter: time.Now().Add(time.Hour * 24 * 366),
    KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
    ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
    BasicConstraintsValid: true,
    IsCA: true,
  }
}

func getClientCertificate() x509.Certificate {
  return x509.Certificate{
    SerialNumber: big.NewInt(1),
    Subject: pkix.Name{
      Organization: []string{"GitLab, Inc."},
      Country: []string{"US"},
    },
    IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
    NotBefore: time.Now(),
    NotAfter: time.Now().Add(time.Hour * 24 * 366),
    KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
    ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
    BasicConstraintsValid: true,
    IsCA: false,
  }
}
