package main

import(
	"log"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
	"bytes"
	"fmt"
	"os"
	"io"
)
func createCA(cn string)(cert *Certificate){
	log.Printf("Creating CA with common Name %s", cn)
	//4096 Bit Keys
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("private key cannot be created: %s", err)
	}
	ca := &x509.Certificate{
		//CA always has the number 1 in our pki
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: cn,
		},
		NotBefore:             time.Now(),
		//10 years valid
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	pub := &priv.PublicKey
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		log.Fatalf("create ca failed", err)
	}

	var request bytes.Buffer
	var privateKey bytes.Buffer
	if err := pem.Encode(&request, &pem.Block{Type: PEMCertificateBlockType, Bytes: ca_b}); err != nil {
		log.Fatalf("Could not read Certificate: %s", err)
	}
	if err := pem.Encode(&privateKey, &pem.Block{Type: PEMRSAPrivateKeyBlockType, Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		log.Fatalf("Could not read Private Key: %s", err)
	}
	log.Println("CA successfully created, not yet written to db")
	return &Certificate{
		UserId: 0,
		CN: cn,
		Type: CertificateTypeCA,
		Serial: ca.SerialNumber.Int64(),
		Private: privateKey.String(),
		Public: request.String(),
		ValidTo: ca.NotAfter,
	}
}
func (ca *CA) createServer(cn string)(retCert *Certificate){
	log.Printf("Creating Server with common Name %s", cn)
	// Get CA private key
	block, _ := pem.Decode([]byte(ca.ca.Private))
	if block == nil {
		log.Fatalf("failed to parse ca private key")
	}

	caKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse ca private key: %s", err)
	}

	caCert, err := ReadCertFromPEM(ca.ca.Public)
	if err != nil {
		log.Fatalf("failed to parse ca cert: %v", err)
	}
	// Create new cert's key
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("private key cannot be created: %s", err)
	}
	// Prepare certificate
	ca.SerialOld = ca.SerialOld + 1
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(ca.SerialOld),
		Subject: pkix.Name{
			CommonName: cn,
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	pub := &priv.PublicKey
	cert_b, err := x509.CreateCertificate(rand.Reader, cert, caCert, pub, caKey)
	if err != nil {
		log.Fatalf("cert not created: %s", err)
	}
	var request bytes.Buffer
	var privateKey bytes.Buffer
	if err := pem.Encode(&request, &pem.Block{Type: PEMCertificateBlockType, Bytes: cert_b}); err != nil {
		log.Fatalf("Could not read Certificate: %s", err)
	}
	if err := pem.Encode(&privateKey, &pem.Block{Type: PEMRSAPrivateKeyBlockType, Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err != nil {
		log.Fatalf("Could not read Private Key: %s", err)
	}
	
	return &Certificate{
		UserId: 0,
		CN: cn,
		Type: CertificateTypeServer,
		Serial: cert.SerialNumber.Int64(),
		Private: privateKey.String(),
		Public: request.String(),
		ValidTo: cert.NotAfter,
	}
}
// ReadCertFromPEM decodes a PEM encoded string into a x509.Certificate.
func ReadCertFromPEM(s string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(s))
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)
	return cert, nil
}
func (crt *Certificate) WriteFileCert() {
	filename := fmt.Sprintf("/docker/data/%s.crt", crt.CN)
	log.Printf("Writing file: %s", filename)
	certOut, err := os.Create(filename)
	if err != nil {
		log.Println(err)
	}
	_, err = io.WriteString(certOut, crt.Public)
	if err != nil {
		log.Println(err)
	}
	certOut.Close()
}
func (key *Certificate) WriteFileKey() {
	filename := fmt.Sprintf("/docker/data/%s.key", key.CN)
	log.Printf("Writing file: %s", filename)
	keyOut, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Println(err)
	}
	_, err = io.WriteString(keyOut, key.Private)
	if err != nil {
		log.Println(err)
	}
	keyOut.Close()
}


