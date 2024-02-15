package main

import(
	"log"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
	"fmt"
	"os"
	"strings"

	"io/ioutil"
)

func createClientForMail(mail string)(retCert *Certificate, err error){
	dbCrt := Certificate{};
	result := getSingleton().dbConn.conn.Where("mail = ?", mail).Find(&dbCrt)
	if (result.Error != nil){
		log.Printf("Datenbank Fehler: %s", result.Error)
		return nil, result.Error
	}
	cn := fmt.Sprintf("%d.%s.%s", result.RowsAffected + 1, mail, os.Getenv("CN_SUFFIX"))
	log.Printf("Creating Client with common Name %s", cn)

	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("private key cannot be created: %s", err)
	}

	// Prepare certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(getSingleton().SerialOld + 1),
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
		Mail: mail,
		CN: cn,
		Type: CertificateTypeClient,
	}
	caCert, caKey := getSingleton().dbConn.loadCA().getKeyAndCertCA()
	retCert.createCert(caCert, cert, priv, caKey)

	result = getSingleton().dbConn.conn.Create(&retCert)
	if(result.Error != nil){
		log.Println(result.Error)
	}else{
		log.Println(result.RowsAffected)
	}

	getSingleton().dbConn.createCCD(mail)

	result = getSingleton().dbConn.conn.Where("mail = ?", mail).Find(&dbCrt)
	if (result.Error != nil){
		log.Printf("Datenbank Fehler: %s", result.Error)
		return nil, result.Error
	}
	log.Printf("%d.%s.%s", result.RowsAffected + 1, mail, os.Getenv("CN_SUFFIX"))

	return retCert, nil
}

func downloadCertByCN(cn string, mail string)(ret []string, err error){
	crt := Certificate{}
	result := getSingleton().dbConn.conn.Where("mail = ? AND cn = ?", mail, cn).First(&crt)
	if (result.Error != nil){
		return nil, fmt.Errorf("Certificate Read error: %s", getSingleton().dbConn.conn.Error)
	}else if (result.RowsAffected == 0){
		return nil, fmt.Errorf("Certification read error (user/certificate not found)")
	}
	
	content, err := ioutil.ReadFile("/docker/data/client.ovpn.proto")
	if err != nil {
		log.Printf("Error reading proto file: %s", err)
	}

	lines := strings.Split(string(content), "\n")
	lines = append(lines, "<ca>")
	lines = append(lines, getSingleton().ca.Public)
	lines = append(lines, "</ca>")

	lines = append(lines, "<key>")
	lines = append(lines, crt.Private)
	lines = append(lines, "</key>")

	lines = append(lines, "<cert>")
	lines = append(lines, crt.Public)
	lines = append(lines, "</cert>")

	return lines, nil
}
