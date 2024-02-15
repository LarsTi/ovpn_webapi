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

func createCA(cn string)(retCert *Certificate){
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

	retCert = &Certificate{
		Mail: "",
		CN: cn,
		Type: CertificateTypeCA,
	}
	
	retCert.createCert(ca, ca, priv, priv)

	//db.conn.Create(&retCert)

	//db.createCCD(mail)

	return retCert
}
func (ca *CA) createClient(user *User)(retCert *Certificate){
	cn := fmt.Sprintf("%d.%s.%s.%s", (ca.SerialOld + 1), user.Surname, user.Name, os.Getenv("CN_SUFFIX"))
	log.Printf("Creating Client with common Name %s", cn)
	// Create new cert's key
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("private key cannot be created: %s", err)
	}
	// Prepare certificate
	
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(ca.SerialOld + 1),
		Subject: pkix.Name{
			CommonName: cn,
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement,
	}
	
	retCert = &Certificate{
		Mail: user.Mail,
		CN: cn,
		Type: CertificateTypeClient,
	}
	caCert, caKey := ca.getKeyAndCertCA()
	retCert.createCert(caCert, cert, priv, caKey)
	
	return retCert
}
func (ca *CA) createServer(cn string)(retCert *Certificate){
	log.Printf("Creating Server with common Name %s", cn)
	// Create new cert's key
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("private key cannot be created: %s", err)
	}
	// Prepare certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(ca.SerialOld + 1),
		Subject: pkix.Name{
			CommonName: cn,
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyAgreement | x509.KeyUsageKeyEncipherment,
	}
	retCert = &Certificate{
		Mail: "",
		CN: cn,
		Type: CertificateTypeServer,
	}
	caCert, caKey := ca.getKeyAndCertCA()
	log.Println("Accessing CA to sign cert with commonName %s\n", retCert.CN)

	retCert.createCert(caCert, cert, priv, caKey)
	
	return retCert
}
func (ca *CA) getKeyAndCertCA() (caCert *x509.Certificate, caKey *rsa.PrivateKey){
	log.Println("Accessing CA Key!")
	// Get CA private key
	block, _ := pem.Decode([]byte(ca.ca.Private))
	if block == nil {
		log.Fatalf("failed to parse ca private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse ca private key: %s", err)
	}

	cert, err := ReadCertFromPEM(ca.ca.Public)
	if err != nil {
		log.Fatalf("failed to parse ca cert: %v", err)
	}
	return cert, key
}
func (crt *Certificate) createCert(caCert, crtCert *x509.Certificate, signedKey, signingKey *rsa.PrivateKey){
	var privateKey bytes.Buffer
	err := pem.Encode(&privateKey, &pem.Block{Type: PEMRSAPrivateKeyBlockType, Bytes: x509.MarshalPKCS1PrivateKey(signedKey)}); 
	if err != nil{
		log.Fatalf("Could not read private Key of common name %s: %s", crt.CN, err)
	}
	crt.Private = privateKey.String()
	crt.Serial = crtCert.SerialNumber.Int64()
	crt.ValidTo = crtCert.NotAfter

	//Signing the cert:
	signedCert, err := x509.CreateCertificate(rand.Reader, crtCert, caCert, &signedKey.PublicKey, signingKey)
	if err != nil {
		log.Fatalf("Could not sign cert with common name %s: %s", crt.CN, err)
	}
	var signed bytes.Buffer
	err = pem.Encode(&signed, &pem.Block{Type: PEMCertificateBlockType, Bytes: signedCert})
	if err != nil {
		log.Fatalf("Could not read signed Key of common name %s: %s", crt.CN, err)
	}
	crt.Public = signed.String()
}

func (ca *CA) createCRL(revoked []Certificate){
	log.Printf("Accessing CA to create CRL\n")
	caCrt, key := ca.getKeyAndCertCA()
	var revokedCertList []pkix.RevokedCertificate
	for _, serial := range revoked {
		revokedCert := pkix.RevokedCertificate{
			SerialNumber:   big.NewInt(serial.Serial),
			RevocationTime: time.Now().UTC(),
		}
		revokedCertList = append(revokedCertList, revokedCert)
	}
	crl, err := caCrt.CreateCRL(rand.Reader, key, revokedCertList, time.Now().UTC(), time.Now().Add(365*24*60*time.Minute).UTC())
	if err != nil {
		log.Printf("CRL: %s", err)
	}
	crlPem := pem.EncodeToMemory(&pem.Block{
		Type:  PEMx509CRLBlockType,
		Bytes: crl,
	})

	filename := fmt.Sprintf("/docker/server/CRL.pem")
	log.Printf("Writing file %s", filename)
	crlOut, err := os.Create(filename)
	if err != nil {
		log.Printf("Write File: %s", err)
	}
	_, err = io.WriteString(crlOut, string(crlPem[:]))
	if err != nil {
		log.Printf("Write String: %s", err)
	}
	crlOut.Close()
}
// ReadCertFromPEM decodes a PEM encoded string into a x509.Certificate.
func ReadCertFromPEM(s string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(s))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("Cert parse: %s", err)
	}
	return cert, nil
}
func (crt *Certificate) WriteFileCert() {
	filename := fmt.Sprintf("/docker/server/%s.crt", crt.CN)
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
	filename := fmt.Sprintf("/docker/server/%s.key", key.CN)
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


