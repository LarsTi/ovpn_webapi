package main

import(
	"gorm.io/gorm"
	"time"
)
const(
	CertificateTypeServer	= 1
	CertificateTypeCA	= 2
	CertificateTypeClient	= 3
	// PEM encoding types

	PEMCertificateBlockType   string = "CERTIFICATE"
	PEMRSAPrivateKeyBlockType        = "RSA PRIVATE KEY"
	PEMx509CRLBlockType              = "X509 CRL"
	PEMCSRBlockType                  = "CERTIFICATE REQUEST"

)
type User struct {
	gorm.Model
	UserId		uint		`json:"user"`
	Name		string		`json:"name"`
	Surname		string		`json:"surname"`
	Org		string		`json:"org"`
	Mail		string		`json:"Mail"`
}
type AccessGroup struct {
	gorm.Model
	Name		string		`json:"name"`
	Subnet		string		`json:"subnet"`
	Mask		string		`json:"mask"`
}
type Certificate struct {
	gorm.Model
	UserId		uint		`json:"user"`		//ref to User.ID
	CN		string		`json:"common_name"`	//Common Name of Cert
	Type		int		`json:"cert_type"`	//one of CertificateType...
	Serial		int64		`json:"-"`		//Serialnumber
	Private		string		`json:"-"`		//Private Key
	Public		string		`json:"public"`		//Public Key
	Revoked		bool		`json:"revoked"`	//true, if it is revoked
	ValidTo		time.Time	`json:"valid_to"`	// Valid to flag
}
type UserAccess struct{
	gorm.Model
	UserId		uint		`json:"user"`		//ref to User.ID
	AccessGroup	uint		`json:"group"`		//ref to AccessGroup.ID
}
type CA struct{
	ca		*Certificate	//CA
	SerialOld	int64		//letzte Vergebene Serial
	db		*DB
}
type DB struct {
	conn		*gorm.DB
}
